//go:build windows

package oscommands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var errConPtyUnavailable = errors.New("ConPTY not available")

// conPty owns a Windows pseudoconsole and the child process attached to it.
//
// Lifecycle:
//  1. newConPty creates two pipes and CreatePseudoConsole, producing hpc.
//  2. Start spawns the child via CreateProcessW with STARTUPINFOEX whose
//     attribute list carries PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE=hpc.
//     The child inherits console handles via the pseudoconsole, NOT via
//     the pipe handles (which we close on our side after CreateProcess).
//  3. Stdin()/Stdout() return our ends of the two pipes.
//  4. Close tears down stdin pipe, ClosePseudoConsole, stdout pipe, and
//     the process/thread handles. Callers must cmd.Wait() before Close.
//
// Concurrency: mu serializes writes to the handle fields against
// snapshot reads in pipeReader/pipeWriter. This provides the Go
// memory-model happens-before relationship needed for safe cross-
// goroutine handle access. In-flight I/O safety is handled by the
// Windows kernel: CloseHandle during an outstanding ReadFile/WriteFile
// aborts the call with ERROR_INVALID_HANDLE or ERROR_BROKEN_PIPE,
// which pipeReader/Writer translate to io.EOF / io.ErrClosedPipe below.
type conPty struct {
	mu          sync.Mutex
	hpc         windows.Handle // pseudoconsole handle
	inRead      windows.Handle // input pipe, child side (closed after spawn)
	inWrite     windows.Handle // input pipe, our side (written to by writer())
	outRead     windows.Handle // output pipe, our side (read from by reader())
	outWrite    windows.Handle // output pipe, child side (closed after spawn)
	attrList    *windows.ProcThreadAttributeListContainer
	processInfo windows.ProcessInformation
	processOpen bool // true once processInfo.Process has been handed to os.Process
}

func newConPty(width, height int16) (*conPty, error) {
	c := &conPty{}

	// Create the two pipes. The child side handles will be consumed by
	// CreatePseudoConsole; our side handles stay with us.
	if err := windows.CreatePipe(&c.inRead, &c.inWrite, nil, 0); err != nil {
		return nil, fmt.Errorf("CreatePipe(in): %w", err)
	}
	if err := windows.CreatePipe(&c.outRead, &c.outWrite, nil, 0); err != nil {
		_ = windows.CloseHandle(c.inRead)
		_ = windows.CloseHandle(c.inWrite)
		return nil, fmt.Errorf("CreatePipe(out): %w", err)
	}

	size := windows.Coord{X: width, Y: height}
	err := windows.CreatePseudoConsole(size, c.inRead, c.outWrite, 0, &c.hpc)
	if err != nil {
		c.closePipesBestEffort()
		// Translate "procedure not found" / "function not implemented" into
		// our sentinel so callers can fall back cleanly.
		if errors.Is(err, windows.ERROR_PROC_NOT_FOUND) ||
			errors.Is(err, windows.ERROR_CALL_NOT_IMPLEMENTED) {
			return nil, errConPtyUnavailable
		}
		return nil, fmt.Errorf("CreatePseudoConsole: %w", err)
	}

	// After CreatePseudoConsole, the child-side handles are duplicated into the
	// pseudoconsole. Close our references so they don't leak; the conpty keeps
	// its own duplicates alive.
	_ = windows.CloseHandle(c.inRead)
	c.inRead = 0
	_ = windows.CloseHandle(c.outWrite)
	c.outWrite = 0

	return c, nil
}

// Start spawns cmd's program attached to this conPty's pseudoconsole.
//
// We bypass cmd.Start() because attaching a process to a pseudoconsole
// requires CreateProcessW with STARTUPINFOEX (not supported by exec.Cmd's
// SysProcAttr). After a successful spawn we populate cmd.Process via
// os.FindProcess so later cmd.Wait() on the caller side works as expected.
func (c *conPty) Start(cmd *exec.Cmd) error {
	if cmd.Process != nil {
		return errors.New("conPty.Start: cmd already started")
	}

	// Self-cleanup on any early return: if we don't reach the successful
	// tail, close everything we've partially set up. Close() is idempotent
	// and guards each field, so it's safe to run on partial state.
	started := false
	defer func() {
		if !started {
			_ = c.Close()
		}
	}()

	// Build the command line. CreateProcess reads argv[0] from the command
	// line, not from lpApplicationName. Prefer cmd.Args when present (this
	// matches stock exec.Cmd behavior on Windows); fall back to cmd.Path
	// for a bare exec.Cmd{}.
	var argv []string
	if len(cmd.Args) > 0 {
		argv = cmd.Args
	} else {
		argv = []string{cmd.Path}
	}
	cmdLine := buildCommandLine(argv)
	cmdLineUTF16, err := syscall.UTF16PtrFromString(cmdLine)
	if err != nil {
		return fmt.Errorf("UTF16PtrFromString(cmdLine): %w", err)
	}
	// We also pass lpApplicationName (argv0) so CreateProcess doesn't have to
	// re-parse it from the command line. Matches the pattern used by
	// aymanbagabas/go-pty and Microsoft's ConPTY sample.
	appNameUTF16, err := syscall.UTF16PtrFromString(cmd.Path)
	if err != nil {
		return fmt.Errorf("UTF16PtrFromString(cmd.Path): %w", err)
	}

	// Environment block. If cmd.Env is nil, use the parent's environment.
	envSlice := cmd.Env
	if envSlice == nil {
		envSlice = os.Environ()
	}
	envBlock, err := buildEnvBlock(envSlice)
	if err != nil {
		return fmt.Errorf("buildEnvBlock: %w", err)
	}

	var cwdPtr *uint16
	if cmd.Dir != "" {
		p, err := syscall.UTF16PtrFromString(cmd.Dir)
		if err != nil {
			return fmt.Errorf("UTF16PtrFromString(Dir): %w", err)
		}
		cwdPtr = p
	}

	// Build a STARTUPINFOEX whose attribute list carries the pseudoconsole.
	attrList, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return fmt.Errorf("NewProcThreadAttributeList: %w", err)
	}
	// PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE expects the HPCON *value* itself as
	// lpValue (not a pointer to a variable holding it). This is the pattern used
	// by Microsoft's ConPTY sample and aymanbagabas/go-pty; passing &c.hpc
	// instead yields an empty child stdout on Windows 11 26100 in practice.
	// Note: this converts a handle value to unsafe.Pointer intentionally —
	// Windows reads sizeof(HPCON) bytes and treats the pointer bits as the
	// handle value, matching the documented ABI.
	if err := attrList.Update(
		windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE,
		unsafe.Pointer(c.hpc),
		unsafe.Sizeof(c.hpc),
	); err != nil {
		attrList.Delete()
		return fmt.Errorf("Update PSEUDOCONSOLE: %w", err)
	}
	c.attrList = attrList

	siEx := new(windows.StartupInfoEx)
	siEx.Flags = windows.STARTF_USESTDHANDLES
	siEx.Cb = uint32(unsafe.Sizeof(*siEx))
	siEx.ProcThreadAttributeList = attrList.List()

	// EXTENDED_STARTUPINFO_PRESENT is required to honor STARTUPINFOEX.
	// CREATE_UNICODE_ENVIRONMENT tells CreateProcess that envBlock is UTF-16.
	flags := uint32(windows.EXTENDED_STARTUPINFO_PRESENT | windows.CREATE_UNICODE_ENVIRONMENT)

	// Passing nil for both procSecurity/threadSecurity yields default security
	// with non-inheritable returned process/thread handles. We don't want lazygit
	// to accidentally leak the child's process/thread handles into future
	// grandchildren it spawns.
	err = windows.CreateProcess(
		appNameUTF16,      // applicationName: explicit argv0
		cmdLineUTF16,      // commandLine
		nil,               // procSecurity: default + non-inheritable
		nil,               // threadSecurity: default + non-inheritable
		false,             // inheritHandles: handles come via the pseudoconsole
		flags,             // creationFlags
		&envBlock[0],      // env block (UTF-16)
		cwdPtr,            // currentDir (may be nil)
		&siEx.StartupInfo, // lpStartupInfo: first field of StartupInfoEx is StartupInfo
		&c.processInfo,    // output
	)
	if err != nil {
		attrList.Delete()
		c.attrList = nil
		return fmt.Errorf("CreateProcess: %w", err)
	}
	// Close thread handle; we only need the process handle.
	_ = windows.CloseHandle(c.processInfo.Thread)
	c.processInfo.Thread = 0

	// Populate cmd.Process so the caller's later cmd.Wait() works. os.FindProcess
	// on Windows opens a fresh handle via OpenProcess — safe to use alongside
	// our own c.processInfo.Process handle. We close ours in Close().
	proc, err := os.FindProcess(int(c.processInfo.ProcessId))
	if err != nil {
		// FindProcess should not fail on a pid we just spawned, but guard anyway.
		_ = windows.TerminateProcess(c.processInfo.Process, 1)
		// Close our handle; without this the kernel keeps the process object
		// alive even after termination, leaking until lazygit exits.
		_ = windows.CloseHandle(c.processInfo.Process)
		c.processInfo.Process = 0
		return fmt.Errorf("os.FindProcess: %w", err)
	}
	cmd.Process = proc
	c.processOpen = true
	started = true

	return nil
}

// Stdin returns an io.Writer for writing to the child's stdin (via the
// pseudoconsole's input pipe).
func (c *conPty) Stdin() io.Writer {
	return pipeWriter{c: c}
}

// Stdout returns an io.Reader for reading the child's stdout+stderr (merged
// via the pseudoconsole's output pipe).
func (c *conPty) Stdout() io.Reader {
	return pipeReader{c: c}
}

// Close tears the conPty down in a well-defined order:
//  1. Close the stdin pipe so the child sees EOF on reads.
//  2. ClosePseudoConsole — signals the child through its attached console;
//     after this, the output pipe writer (held by the console host) closes,
//     which unblocks our reader goroutine with EOF.
//  3. Close the stdout pipe on our side.
//  4. Close the attribute list.
//  5. Close our copy of the process handle (caller's cmd.Wait uses its own).
//
// The caller is expected to have called cmd.Wait() before Close; this
// Close only tears down handles and the pseudoconsole — it does not wait
// for the child. Close is idempotent and safe to call concurrently with
// in-flight Read/Write on Stdin()/Stdout() (the pipe handle protects
// readers/writers via c.mu).
func (c *conPty) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var errs []error
	if c.inWrite != 0 {
		if err := windows.CloseHandle(c.inWrite); err != nil {
			errs = append(errs, fmt.Errorf("close inWrite: %w", err))
		}
		c.inWrite = 0
	}
	if c.hpc != 0 {
		// ClosePseudoConsole is void-returning in the Windows API.
		windows.ClosePseudoConsole(c.hpc)
		c.hpc = 0
	}
	if c.outRead != 0 {
		if err := windows.CloseHandle(c.outRead); err != nil {
			errs = append(errs, fmt.Errorf("close outRead: %w", err))
		}
		c.outRead = 0
	}
	if c.attrList != nil {
		// Delete is void-returning.
		c.attrList.Delete()
		c.attrList = nil
	}
	if c.processOpen && c.processInfo.Process != 0 {
		if err := windows.CloseHandle(c.processInfo.Process); err != nil {
			errs = append(errs, fmt.Errorf("close process: %w", err))
		}
		c.processInfo.Process = 0
		c.processOpen = false
	}
	return errors.Join(errs...)
}

func (c *conPty) closePipesBestEffort() {
	for _, h := range [...]*windows.Handle{&c.inRead, &c.inWrite, &c.outRead, &c.outWrite} {
		if *h != 0 {
			_ = windows.CloseHandle(*h)
			*h = 0
		}
	}
}

// pipeWriter / pipeReader wrap a windows.Handle as io.Writer / io.Reader.
// They use windows.WriteFile / windows.ReadFile directly (rather than
// os.NewFile) so we don't double-close the handle when conPty.Close runs.
//
// Both hold a reference to the parent conPty and snapshot the handle
// under c.mu before issuing the syscall. The mutex provides the Go
// memory-model synchronization for the handle-field read; it does NOT
// attempt to prevent in-flight handle reuse. That safety comes from
// the Windows kernel, which keeps the underlying object alive while
// any ReadFile/WriteFile is still outstanding.

type pipeWriter struct{ c *conPty }

func (w pipeWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	w.c.mu.Lock()
	h := w.c.inWrite
	w.c.mu.Unlock()
	if h == 0 {
		return 0, io.ErrClosedPipe
	}
	var n uint32
	if err := windows.WriteFile(h, p, &n, nil); err != nil {
		// If Close() fired between our snapshot and the syscall, Windows
		// returns ERROR_INVALID_HANDLE. Surface that as ErrClosedPipe for
		// callers; ERROR_BROKEN_PIPE means the reader (child) exited.
		if errors.Is(err, windows.ERROR_INVALID_HANDLE) ||
			errors.Is(err, windows.ERROR_BROKEN_PIPE) {
			return int(n), io.ErrClosedPipe
		}
		return int(n), err
	}
	return int(n), nil
}

type pipeReader struct{ c *conPty }

func (r pipeReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	// Snapshot the handle under c.mu — Close() may zero c.outRead concurrently.
	r.c.mu.Lock()
	h := r.c.outRead
	r.c.mu.Unlock()
	if h == 0 {
		return 0, io.EOF
	}
	var n uint32
	err := windows.ReadFile(h, p, &n, nil)
	if err != nil {
		// When the pseudoconsole is closed, the peer writer closes and ReadFile
		// returns ERROR_BROKEN_PIPE. When Close() races ahead and zeros our
		// handle after we read it, ReadFile returns ERROR_INVALID_HANDLE.
		// Translate both to EOF for Go consumers.
		if errors.Is(err, windows.ERROR_BROKEN_PIPE) ||
			errors.Is(err, windows.ERROR_INVALID_HANDLE) {
			return int(n), io.EOF
		}
		return int(n), err
	}
	if n == 0 {
		return 0, io.EOF
	}
	return int(n), nil
}

// buildCommandLine joins argv with spaces, quoting each element as needed.
func buildCommandLine(argv []string) string {
	parts := make([]string, len(argv))
	for i, a := range argv {
		parts[i] = syscall.EscapeArg(a)
	}
	return strings.Join(parts, " ")
}

// buildEnvBlock converts a []string of "KEY=VALUE" entries into a UTF-16
// null-separated, double-null-terminated environment block as required by
// CreateProcess with CREATE_UNICODE_ENVIRONMENT.
func buildEnvBlock(env []string) ([]uint16, error) {
	if len(env) == 0 {
		// Empty environment block is still a double-NUL.
		return []uint16{0, 0}, nil
	}
	// Rough worst-case sizing; grown as needed.
	out := make([]uint16, 0, 1024)
	for _, e := range env {
		u16, err := syscall.UTF16FromString(e)
		if err != nil {
			return nil, err
		}
		// UTF16FromString appends a terminating NUL; keep it — that's the separator.
		out = append(out, u16...)
	}
	// Final NUL turns the trailing separator into a double-NUL terminator.
	out = append(out, 0)
	return out, nil
}

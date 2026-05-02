//go:build windows

package oscommands

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/mgutz/str"
)

// /s before /c makes cmd.exe always strip exactly the first and last
// quote on the line.
// The EscapeArg-built inner quoting around each arg survives.
// The str.ToArgv → syscall.EscapeArg round-trip is intentional, not
// redundant: lazygit's Quote uses `\"` (bash-compatible quoting),
// which str.ToArgv has explicit Windows handling for (it drops the `\`
// before `"`).
// EscapeArg then re-quotes the clean args with `"` for
// CommandLineToArgvW.
func (self *CmdObjBuilder) newShellCmd(payload string) *CmdObj {
	args := str.ToArgv(fmt.Sprintf("%s %s %s", self.platform.Shell, self.platform.ShellArg, payload))
	cmdObj := self.New(args)
	cmdObj.cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: buildShellCmdLine(args),
	}
	return cmdObj
}

// Assembles `<shell> /s <shellArg> "<inner>"` where <inner> is
// args[2:] joined with spaces, each escaped via EscapeArg.
func buildShellCmdLine(args []string) string {
	var inner strings.Builder
	for i, a := range args[2:] {
		if i > 0 {
			inner.WriteByte(' ')
		}
		inner.WriteString(syscall.EscapeArg(a))
	}
	return syscall.EscapeArg(args[0]) + " /s " + syscall.EscapeArg(args[1]) +
		` "` + inner.String() + `"`
}

//go:build windows

package oscommands

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Spawns `cmd.exe /c echo hello` under a ConPTY, reads its output through
// the ConPTY's stdout pipe, and asserts we see "hello" in the output.
func TestConPty_EchoRoundTrip(t *testing.T) {
	p, err := newConPty(80, 24)
	if err != nil {
		if errors.Is(err, errConPtyUnavailable) {
			t.Skip("ConPTY not available on this Windows version")
		}
		t.Fatalf("newConPty: %v", err)
	}
	shell := os.Getenv("ComSpec")
	if shell == "" {
		shell = "cmd.exe"
	}
	cmd := exec.Command(shell, "/c", "echo", "hello")
	if err := p.Start(cmd); err != nil {
		_ = p.Close()
		t.Fatalf("Start: %v", err)
	}

	outCh := make(chan []byte, 1)
	go func() {
		b, _ := io.ReadAll(p.Stdout())
		outCh <- b
	}()

	waitErr := cmd.Wait()
	closeErr := p.Close()

	select {
	case b := <-outCh:
		if !strings.Contains(string(b), "hello") {
			t.Errorf("stdout did not contain \"hello\": %q", string(b))
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout draining stdout after Close")
	}

	if waitErr != nil {
		t.Errorf("cmd.Wait: %v", waitErr)
	}
	if closeErr != nil {
		t.Errorf("Close: %v", closeErr)
	}
}

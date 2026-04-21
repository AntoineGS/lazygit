package oscommands

import (
	"errors"
	"fmt"
	"os/exec"
)

// defaultConPtyCols / defaultConPtyRows sized for credential prompts.
// Lazygit does not currently forward terminal resize events to child
// processes on any platform, so a fixed default is fine here too.
const (
	defaultConPtyCols int16 = 80
	defaultConPtyRows int16 = 24
)

// getCmdHandlerPty runs cmd under a Windows pseudoconsole (ConPTY) so that
// interactive programs (git, ssh) see a TTY and write credential prompts
// to stdout where runAndDetectCredentialRequest can parse them.
//
// If ConPTY is unavailable (pre-Windows-10-1809, or a locked-down host),
// we fall back to getCmdHandlerNonPty. That matches today's behavior —
// the command runs but no credential popup appears. Users can still see
// the stuck state via the command log; no regression vs the previous stub.
func (self *cmdObjRunner) getCmdHandlerPty(cmd *exec.Cmd) (*cmdHandler, error) {
	p, err := newConPty(defaultConPtyCols, defaultConPtyRows)
	if err != nil {
		if errors.Is(err, errConPtyUnavailable) {
			self.log.WithError(err).Warn("ConPTY unavailable, falling back to non-PTY handler")
			return self.getCmdHandlerNonPty(cmd)
		}
		return nil, err
	}
	if err := p.Start(cmd); err != nil {
		_ = p.Close()
		return nil, fmt.Errorf("conpty start: %w", err)
	}
	return &cmdHandler{
		stdoutPipe: newAnsiStripReader(p.Stdout()),
		stdinPipe:  p.Stdin(),
		close:      p.Close,
	}, nil
}

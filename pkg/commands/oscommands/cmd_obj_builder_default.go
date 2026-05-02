//go:build !windows

package oscommands

import (
	"fmt"

	"github.com/mgutz/str"
)

func (self *CmdObjBuilder) newShellCmd(payload string) *CmdObj {
	args := str.ToArgv(fmt.Sprintf("%s %s %s", self.platform.Shell, self.platform.ShellArg, payload))
	return self.New(args)
}

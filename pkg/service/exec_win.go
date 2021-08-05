// +build windows

package service

import (
	"os/exec"
)

func (s *Controller) exec(args ...string) error {
	cmd := exec.Command("explorer.exe", args...)

	cmd.Dir = s.basedir

	dbg(s.basedir)

	return cmd.Start()
}

// +build !windows

package service

import (
	"os/exec"
)

// exec is not implemented on non windows systems
func (s *Controller) exec(args ...string) error {
	return nil
}

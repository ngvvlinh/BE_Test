// +build !windows

package service

// FileOpen is not implemented on non windows systems
func (s *Controller) FileOpen(name string) error {
	return nil
}

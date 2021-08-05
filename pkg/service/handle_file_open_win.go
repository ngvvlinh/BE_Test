// +build windows

package service

import (
	"os"
	"strings"
)

func (s *Controller) FileOpen(name string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)

	file, err := s.tempFs.OpenFile(_name, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	defer file.Close()

	return s.exec(file.Name())
}

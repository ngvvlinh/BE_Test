package handler

import (
	"os"

	"dynexo.de/pkg/log"

	"github.com/spf13/afero"
)

type (
	Dir struct {
		fs       afero.Fs
		tempFs   afero.Fs
		readOnly bool
		logger   log.ILogger
	}
)

func NewDirHandler(fs afero.Fs, tempFs afero.Fs, readOnly bool, logger log.ILogger) (*Dir, error) {

	var err error

	s := &Dir{
		fs:       fs,
		tempFs:   tempFs,
		readOnly: readOnly,
		logger:   logger,
	}

	return s, err
}

func (s *Dir) Create(name string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.fs.MkdirAll(name, 0777)
}

func (s *Dir) Remove(name string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.fs.RemoveAll(name)
}

func (s *Dir) Stat(name string) ([]os.FileInfo, error) {

	file, err := s.fs.OpenFile(name, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	return file.Readdir(0)
}

func (s *Dir) Readdir(name string) ([]os.FileInfo, error) {

	return s.Stat(name)
}

func (s *Dir) Duplicate(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	// TODO

	return nil
}

func (s *Dir) Copy(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	// TODO

	return nil
}

func (s *Dir) Move(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.fs.Rename(name, name2)

}

func (s *Dir) Rename(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.fs.Rename(name, name2)
}

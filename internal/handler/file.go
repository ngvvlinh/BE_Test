package handler

import (
	"fmt"
	"os"
	"time"

	dio "dynexo.de/pkg/io"
	"dynexo.de/pkg/log"

	"github.com/kr/pretty"
	"github.com/spf13/afero"
)

var (
	dbg = pretty.Println
)

type (
	File struct {
		fs       afero.Fs
		tempFs   afero.Fs
		readOnly bool
		logger   log.ILogger
	}
)

func NewFileHandler(fs afero.Fs, tempFs afero.Fs, readOnly bool, logger log.ILogger) (*File, error) {

	var err error

	s := &File{
		fs:       fs,
		tempFs:   tempFs,
		readOnly: readOnly,
		logger:   logger,
	}

	return s, err
}

func (s *File) Create(name string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	_, err := s.fs.Create(name)

	return err
}

func (s *File) Remove(name string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.fs.Remove(name)
}

func (s *File) Stat(name string) (os.FileInfo, error) {

	return s.fs.Stat(name)
}

func (s *File) rename(name string, name2 string) error {
	return s.fs.Rename(name, name2)
}

func (s *File) copy(name string, name2 string) error {

	file1, err1 := s.fs.OpenFile(name, os.O_RDONLY, 0666)
	if err1 != nil {
		return err1
	}

	file2, err2 := s.fs.OpenFile(name2, os.O_TRUNC|os.O_RDWR, 0666)
	if err2 != nil {
		return err2
	}

	_, err := dio.CopyBuffer(file2, file1, nil)

	return err
}

func (s *File) Duplicate(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.copy(name, name2)
}

func (s *File) Copy(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.copy(name, name2)
}

func (s *File) Move(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.rename(name, name2)
}

func (s *File) Rename(name string, name2 string) error {

	if s.readOnly {
		return os.ErrPermission
	}

	return s.rename(name, name2)
}

func (s *File) Lock(name string, user string) error {

	_name := fmt.Sprintf("%s.ulock", name)

	if st, err := s.fs.Stat(_name); err == nil && st.Size() > 0 {
		return os.ErrExist
	}

	if file, err := s.fs.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666); err == nil {
		str := fmt.Sprintf("{\"u\": \"%s\", \"t\": %d}\n", user, time.Now().UnixNano())
		file.WriteString(str)
		file.Sync()
		file.Close()
	} else {
		return err
	}

	return nil
}

func (s *File) Unlock(name string) error {

	_name := fmt.Sprintf("%s.ulock", name)

	if st, err := s.fs.Stat(_name); err != nil || st.Size() < 1 {
		return os.ErrNotExist
	}

	return s.Remove(name)
}

func (s *File) Share(name string) (string, error) {

	return "", nil
}

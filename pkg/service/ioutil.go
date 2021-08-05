package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"dynexo.de/pkg/io"

	"github.com/spf13/afero"
)

var (
	ErrSizeLimit = errors.New("Size limit exceeded")
	ErrNoEntry   = errors.New("Entry is missing")
)

type (
	FileInfo struct {
		name    string
		size    int64
		mode    os.FileMode
		modtime time.Time
	}
)

func (fi FileInfo) Name() string       { return fi.name }
func (fi FileInfo) Size() int64        { return fi.size }
func (fi FileInfo) Mode() os.FileMode  { return fi.mode }
func (fi FileInfo) ModTime() time.Time { return fi.modtime }
func (fi FileInfo) IsDir() bool        { return fi.mode.IsDir() }
func (fi FileInfo) Sys() interface{}   { return nil }

func (fi FileInfo) String() string {
	str := `{"name":"%s", "size":%d, "mtime":%d, "dir":%v}`
	return fmt.Sprintf(str, fi.Name(), fi.Size(), fi.modtime.UnixNano(), fi.IsDir())
}

//
// fileCopy
//

func (s *Controller) fileCopy(srcFs, dstFs afero.Fs, name, name2 string) error {

	file1, er1 := srcFs.OpenFile(name, os.O_RDONLY, 0666)
	if er1 != nil {
		return er1
	}

	defer file1.Close()

	file2, er2 := dstFs.OpenFile(name2, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if er2 != nil {
		return er2
	}

	defer file2.Close()

	buf32 := s.bpool.Get32()
	defer s.bpool.Put32(buf32)

	_, err := io.CopyBuffer(file2, file1, buf32)
	if err != nil {
		return err
	}

	return err
}

package io

import (
	"fmt"
	"os"
	"time"
)

type (
	FileInfo struct {
		FileName    string
		FileSize    int64
		FileMode    os.FileMode
		FileModtime time.Time
	}
)

func (fi FileInfo) Name() string       { return fi.FileName }
func (fi FileInfo) Size() int64        { return fi.FileSize }
func (fi FileInfo) Mode() os.FileMode  { return fi.FileMode }
func (fi FileInfo) ModTime() time.Time { return fi.FileModtime }
func (fi FileInfo) IsDir() bool        { return fi.FileMode.IsDir() }
func (fi FileInfo) Sys() interface{}   { return nil }

func (fi FileInfo) String() string {
	str := `{"name":"%s", "size":%d, "mtime":%d, "dir":%v}`
	return fmt.Sprintf(str, fi.Name(), fi.Size(), fi.FileModtime.UnixNano(), fi.IsDir())
}

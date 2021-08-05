package io

import (
	"os"
)

type FileInfoByName []os.FileInfo

func (f FileInfoByName) Len() int           { return len(f) }
func (f FileInfoByName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f FileInfoByName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

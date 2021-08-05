package io

import (
	"os"

	"dynexo.de/pkg/io"

	"github.com/spf13/afero"
)

func FileCopy(srcFs, dstFs afero.Fs, name, name2 string, buf []byte) error {

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

	// CopyBuffer allocates mem if buf is nil

	_, err := io.CopyBuffer(file2, file1, buf)
	if err != nil {
		return err
	}

	return err
}

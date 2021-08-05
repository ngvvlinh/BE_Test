package umodel

import (
	"dynexo.de/pkg/log"
	"github.com/spf13/afero"
)

type (
	IController interface {
		Debug()
		Serve()
		Fs() afero.Fs
		FsRead() afero.Fs
		FsWrite() afero.Fs
		Log() log.ILogger
		Messages() []string

		HandleFile([]byte) (IRes, error)
		HandleFolder([]byte) (IRes, error)
		HandleSearch([]byte) (IRes, error)
		HandleRecenlyUsed([]byte) (IRes, error)
	}
)

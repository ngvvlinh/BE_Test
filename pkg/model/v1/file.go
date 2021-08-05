package umodel

// import (
// 	"dynexo.de/pkg/log"
// 	"dynexo.de/ufyle/pkg/err"
// )

type (
	File interface {
		FileCreate(string) error
		FileRemove(string) error
		FileStat(string) ([]byte, error)
		FileDuplicate(string, string) error
		FileCopy(string, string) error
		FileMove(string, string) error
		FileRename(string, string) error

		FileLock(string, string) error
		FileUnlock(string) error
		FileShare(string) (string, error)
	}
)

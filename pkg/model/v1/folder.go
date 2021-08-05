package umodel

// import (
// 	"dynexo.de/pkg/log"
// 	"dynexo.de/ufyle/pkg/err"
// )

type (
	Folder interface {
		DirCreate(string) error
		DirRemove(string) error
		DirStat(string) ([]byte, error)
		DirReaddir(string) ([]byte, error)
		DirDuplicate(string, string) error
		DirCopy(string, string) error
		DirMove(string, string) error
		DirRename(string, string) error
	}
)

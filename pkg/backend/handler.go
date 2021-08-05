package backend

import (
	uerr "dynexo.de/ufyle/pkg/err"
	"dynexo.de/ufyle/pkg/model/v1"
	//"dynexo.de/ufyle/pkg/handler/file"
	//"dynexo.de/ufyle/pkg/handler/folder"
)

type (
	HandlerFunc func(umodel.Data) ([]byte, error)
)

// HandleFile

const (
	OpFileStat int = iota
	OpFileCreate
	OpFileRemove
	OpFileReserved1
	OpFileDuplicate
	OpFileCopy
	OpFileMove
	OpFileRename
	OpFileRun
	OpFileUpload
	OpFileDownload

	OpFileLock int = iota + 12
	OpFileUnlock

	OpFileShare int = iota + 16
)

// func (s *Controller) HandleFile(data []byte) (umodel.HRes, error) {

// 	s.logger.Debug("HandleFile:", string(data))

// 	udata, er := umodel.UnmarshalData(data)
// 	if er != nil {
// 		s.logger.Error("HandleFile:", er)
// 		return nil, er
// 	}

// 	// 0	create
// 	// 1 	remove			recursively remove all files and folders based on defined path
// 	// 2 	stat			retrieve folder stats (date, size, last changed); returns error if not existing
// 	// 3	reserved
// 	// 4 	duplicate		duplicate file
// 	// 5 	copy			copy file
// 	// 6 	move			move file
// 	// 7 	rename			rename file
// 	// 8	download		download file
// 	// 9	run				get file and start with OS file handler
// 	// 10
// 	// 11
// 	// 12	lock			lock file for modifications by others
// 	// 13	unlock			unlock file
// 	// 14
// 	// 15
// 	// 16	share			create share link

// 	var err error
// 	var aresp []map[string]interface{}
// 	var resp map[string]interface{}

// 	for i := range udata {

// 		switch udata[i].T {
// 		case OpFileCreate:
// 			err = s.FileCreate(udata[i].Val1)

// 		case OpFileRemove:
// 			err = s.FileRemove(udata[i].Val1)

// 		case OpFileStat:
// 			resp, err = s.FileStat(udata[i].Val1)
// 			if err != nil && resp != nil {
// 				resp["d"] = udata[i]
// 				aresp = append(aresp, resp)
// 			}

// 		case OpFileDuplicate:
// 			err = s.FileDuplicate(udata[i].Val1, udata[i].Val2)

// 		case OpFileMove:
// 			err = s.FileMove(udata[i].Val1, udata[i].Val2)

// 		case OpFileCopy:
// 			err = s.FileCopy(udata[i].Val1, udata[i].Val2)

// 		case OpFileRename:
// 			err = s.FileRename(udata[i].Val1, udata[i].Val2)

// 		case OpFileDownload:
// 			err = s.FileDownload(udata[i].Val1)

// 		}

// 	}

// 	if aresp != nil && len(aresp) < 2 {
// 		delete(aresp[0], "d")
// 	}

// 	dbg(err)

// 	dbg(udata)

// 	return aresp, err
// }

func (s *Controller) HandleFile(data []byte) (umodel.IRes, error) {

	s.logger.Debug("HandleFile:", string(data))

	udata, er := umodel.UnmarshalData(data)
	if er != nil {
		s.logger.Error("HandleFile:", er)
		return nil, er
	}

	// 0	create
	// 1 	remove			recursively remove all files and folders based on defined path
	// 2 	stat			retrieve folder stats (date, size, last changed); returns error if not existing
	// 3	reserved
	// 4 	duplicate		duplicate file
	// 5 	copy			copy file
	// 6 	move			move file
	// 7 	rename			rename file
	// 8	download		download file
	// 9	run				get file and start with OS file handler
	// 10
	// 11
	// 12	lock			lock file for modifications by others
	// 13	unlock			unlock file
	// 14
	// 15
	// 16	share			create share link

	var err error
	var resp map[string]interface{}

	switch udata.T {
	case OpFileCreate:
		err = s.FileCreate(udata.Val1)

	case OpFileRemove:
		err = s.FileRemove(udata.Val1)

	case OpFileStat:
		resp, err = s.FileStat(udata.Val1)

	case OpFileDuplicate:
		err = s.FileDuplicate(udata.Val1, udata.Val2)

	case OpFileMove:
		err = s.FileMove(udata.Val1, udata.Val2)

	case OpFileCopy:
		err = s.FileCopy(udata.Val1, udata.Val2)

	case OpFileRename:
		err = s.FileRename(udata.Val1, udata.Val2)

	case OpFileDownload:
		//err = s.FileDownload(udata.Val1)
	}

	dbg(err)
	dbg(udata)

	aresp := []interface{}{
		resp,
	}

	return aresp, err
}

//
// HandleFolder
//

const (
	OpDirStat int = iota
	OpDirCreate
	OpDirRemove
	OpDirReaddir
	OpDirDuplicate
	OpDirCopy
	OpDirMove
	OpDirRename
)

func (s *Controller) HandleFolder(data []byte) (umodel.IRes, error) {

	s.logger.Debug("HandleFolder:", string(data))

	udata, er := umodel.UnmarshalData(data)
	if er != nil {
		s.logger.Error("HandleFolder:", er)
		return nil, er
	}

	// 0 	mkdir			create a new directory
	// 1 	remove			recursively remove all files and folders based on defined path
	// 2 	stat			retrieve folder stats (date, size, last changed); returns error if not existing
	// 3	readdir			returns a list of fileinfo for all elements of the folder
	// 4 	duplicate		copy folder
	// 5 	copy			copy folder
	// 6 	move			move folder
	// 7 	rename			rename folder

	var err error
	var resp map[string]interface{}

	switch udata.T {
	case OpDirCreate:
		err = s.DirCreate(udata.Val1)

	case OpDirRemove:
		err = s.DirRemove(udata.Val1)

	case OpDirStat:
		resp, err = s.DirStat(udata.Val1)

	case OpDirReaddir:
		err = s.FileDuplicate(udata.Val1, udata.Val2)

	case OpDirDuplicate:
		err = s.DirDuplicate(udata.Val1, udata.Val2)

	case OpDirCopy:
		err = s.DirCopy(udata.Val1, udata.Val2)

	case OpDirRename:
		err = s.DirRename(udata.Val1, udata.Val2)
	}

	dbg(err)
	dbg(udata)

	aresp := []interface{}{
		resp,
	}

	return aresp, err
}

//
// HandleSearch
//

const (
	OpFindGlob int = iota
	OpFindRegex
)

func (s *Controller) HandleSearch(data []byte) (umodel.IRes, error) {

	s.logger.Debug("HandleSearch:", string(data))

	udata, er := umodel.UnmarshalData(data)
	if er != nil {
		s.logger.Error("HandleSearch:", er)
		return nil, er
	}

	var err error
	var resp []string

	if len(data) < 1 {
		s.logger.Error("HandleSearch: Empty request")
		return nil, uerr.ErrNoEntry
	}

	switch udata.T {
	case OpFindGlob:
		resp, err = s.FindGlob(udata.Val1, udata.Val2, udata.Num1)

	case OpFindRegex:
		resp, err = s.FindRegex(udata.Val1, udata.Val2, udata.Num1)
	}

	dbg(err)
	dbg(udata)

	aresp := []interface{}{
		resp,
	}

	return aresp, err
}

//
// HandleRecenlyUsed
//

func (s *Controller) HandleRecenlyUsed(data []byte) (umodel.IRes, error) {

	s.logger.Debug("HandleRecenlyUsed:", string(data))

	return nil, nil
}

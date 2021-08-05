package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"dynexo.de/ufyle/pkg/vars"
)

var (
	files_ignore = []string{
		vars.File_Lock_Suffix,
	}

	folders_ignore = []string{}
)

func formatPath(name string) string {
	return strings.ReplaceAll(name, "/", string(os.PathSeparator))
}

func (s *Controller) FileCreate(name string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)

	path := filepath.Dir(_name)

	if _, err := s.tempFs.Stat(path); err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ERROR_FILE_NOT_FOUND {
			s.tempFs.MkdirAll(path, 0777)
		}
	}

	file, err := s.tempFs.Create(_name)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return s.FileOpen(name)
}

func (s *Controller) FileStat(name string) (map[string]interface{}, error) {

	if strings.HasSuffix(name, file_lock_suffix) {
		return nil, os.ErrPermission
	}

	_name := formatPath(name)

	st1, er1 := s.tempFs.Stat(_name)
	if er1 != nil {
		return nil, er1
	}

	// NOTE stat makes more sense on the real/remote fs instead of the local one

	st, err := s.fs.Stat(_name)
	if err != nil {
		return nil, err
	}

	// NOTE upload from TempFs to Fs in the background
	if st.ModTime().Before(st1.ModTime()) {

		s.AddMessage(fmt.Sprintf("Uploading temporary file %s", _name))

		go s.fileCopy(s.tempFs, s.fs, _name, _name)
	}

	data := map[string]interface{}{
		"name": st.Name(),
		"size": st.Size(),
		"time": st.ModTime().UnixNano(),
		"dir":  st.IsDir(),
	}

	return data, err
}

func (s *Controller) FileRemove(name string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	var err error

	_name := formatPath(name)

	// NOTE Tempfs

	if err = s.tempFs.Remove(_name); err != nil {
		return err
	}

	// NOTE Fs

	if err = s.fs.Remove(_name); err != nil {
		return err
	}

	return err
}

func (s *Controller) FileDuplicate(name, name2 string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	if strings.HasSuffix(name2, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)
	_name2 := formatPath(name2)

	if len(_name2) < 1 {
		return filepath.ErrBadPattern
	}

	go s.fileCopy(s.tempFs, s.tempFs, _name, _name2)

	return s.fileCopy(s.fs, s.fs, _name, _name2)
}

func (s *Controller) FileCopy(name, name2 string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	if strings.HasSuffix(name2, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)
	_name2 := formatPath(name2)

	if len(_name2) < 1 {
		return filepath.ErrBadPattern
	}

	go s.fileCopy(s.tempFs, s.tempFs, _name, _name2)

	return s.fileCopy(s.fs, s.fs, _name, _name2)
}

func (s *Controller) FileMove(name, name2 string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	if strings.HasSuffix(name2, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)
	_name2 := formatPath(name2)

	if len(_name2) < 1 {
		return filepath.ErrBadPattern
	}

	if err := s.tempFs.Rename(_name, _name2); err != nil {
		return err
	}

	if err := s.fs.Rename(_name, _name2); err != nil {
		return err
	}

	return nil
}

func (s *Controller) FileRename(name, name2 string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	if strings.HasSuffix(name2, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)
	_name2 := formatPath(name2)

	if len(_name2) < 1 {
		return filepath.ErrBadPattern
	}

	if err := s.tempFs.Rename(_name, _name2); err != nil {
		return err
	}

	if err := s.fs.Rename(_name, _name2); err != nil {
		return err
	}

	return nil
}

func (s *Controller) FileDownload(name string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	_name := formatPath(name)

	return s.fileCopy(s.fs, s.tempFs, _name, _name)
}

func (s *Controller) FileRun(name string) error {

	if strings.HasSuffix(name, file_lock_suffix) {
		return os.ErrPermission
	}

	if err := s.FileDownload(name); err != nil {
		return err
	}

	return s.FileOpen(name)
}

func (s *Controller) FileLock(name, user string) error {

	_name := fmt.Sprintf("%s%s", formatPath(name), file_lock_suffix)

	if _, err := s.fs.Stat(_name); err == nil {
		return os.ErrExist
	}

	ts := time.Now().UnixNano()

	file, err := s.fs.OpenFile(_name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	file.WriteString(fmt.Sprintf("%d,%s", ts, user))

	return file.Close()
}

func (s *Controller) FileUnlock(name string) error {

	_name := fmt.Sprintf("%s%s", formatPath(name), file_lock_suffix)

	if _, err := s.fs.Stat(_name); err != nil {
		return os.ErrNotExist
	}

	return s.fs.Remove(_name)
}

func (s *Controller) FileShare(name string) (map[string]interface{}, error) {

	// TODO tell remote server to create single use GUID return

	data := map[string]interface{}{
		"link":  "WXfv57nRaAWMvpd04sbiEq2nIgrWCjGU",
		"file":  name,
		"tmout": time.Now().Add(72 * time.Hour).UnixNano(),
	}

	return data, nil
}

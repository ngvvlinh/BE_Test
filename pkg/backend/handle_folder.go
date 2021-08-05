package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	uerr "dynexo.de/ufyle/pkg/err"
	uio "dynexo.de/ufyle/pkg/io"

	"github.com/spf13/afero"
)

func (s *Controller) DirCreate(name string) error {

	_name := formatPath(name)

	path := filepath.Dir(_name)

	if _, err := s.fs.Stat(path); err != nil {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ERROR_FILE_NOT_FOUND {
			s.fs.MkdirAll(path, 0777)
		}
	}

	return nil
}

func (s *Controller) DirRemove(name string) error {

	var err error

	_name := formatPath(name)

	if err = s.fs.Remove(_name); err != nil {
		return err
	}

	return err
}

func (s *Controller) DirStat(name string) (map[string]interface{}, error) {

	_name := formatPath(name)

	st, err := s.fs.Stat(_name)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		"size": st.Size(),
		"time": st.ModTime().UnixNano(),
		"dir":  st.IsDir(),
	}

	return data, err
}

func (s *Controller) readddir(name string, limit int) ([]os.FileInfo, error) {

	_name := formatPath(name)

	f, err := s.fs.Open(_name)
	if err != nil {
		return nil, err
	}

	st, err := f.Readdir(limit)
	f.Close()
	if err != nil {
		return nil, err
	}

	num := len(st)

	for i := 0; i < num; i++ {

		// match against ignored files

		if !st[i].IsDir() {

			for _, match := range files_ignore {
				if strings.HasSuffix(st[i].Name(), match) {
					st[i] = st[num-1]
					num--
					break
				}
			}

		} else {

			// match against ignored folders
			for _, match := range folders_ignore {
				if strings.HasSuffix(st[i].Name(), match) {
					st[i] = st[num-1]
					num--
					break
				}
			}
		}
	}

	st = st[:num]
	sort.Sort(uio.FileInfoByName(st))

	return st, err
}

func (s *Controller) walkdir(name string, limit int) ([]uio.FileInfo, error) {

	var err error
	var st []uio.FileInfo

	_name := formatPath(name)

	if limit < 1 {
		limit = max_files
	}

	st = make([]uio.FileInfo, limit)
	i := 0

	walker := func(path string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		st[i] = uio.FileInfo{
			//name:    fi.Name(),
			FileName:    path,
			FileSize:    fi.Size(),
			FileModtime: fi.ModTime(),
			FileMode:    fi.Mode(),
		}

		i++

		if i >= limit {
			i--
			return uerr.ErrSizeLimit
		}

		return nil
	}

	if err = afero.Walk(s.fs, _name, walker); err != nil && err != uerr.ErrSizeLimit {
		return nil, err
	}

	return st[:i], nil
}

func (s *Controller) copier(srcFs, dstFs afero.Fs, name, name2 string) (int, []*os.PathError, error) {

	// TODO future release: change to limited parallel copier

	var err error

	_name := strings.TrimPrefix(formatPath(name), ".")
	_name2 := strings.TrimPrefix(formatPath(name2), ".")

	var errs []*os.PathError

	cnt := 0

	buf32 := s.bpool.Get32()
	defer s.bpool.Put32(buf32)

	walker := func(path string, fi os.FileInfo, err error) error {

		if err != nil {
			errs = append(errs, &os.PathError{"stat", path, err})
			return nil
		}

		if path == "." {
			return nil
		}

		_path := fmt.Sprintf("%s%s", _name2, strings.TrimPrefix(path, _name))

		if len(_path) < 1 || len(path) < 1 {
			return nil
		}

		//dbg(fi.IsDir(), path, _path)

		if fi.IsDir() {

			if err := dstFs.MkdirAll(path, 0777); err != nil {
				errs = append(errs, &os.PathError{"mkdir", path, err})
			}

		} else {

			if err := uio.FileCopy(srcFs, dstFs, path, path, buf32); err != nil {
				errs = append(errs, &os.PathError{"copy", path, err})
			} else {
				cnt++
			}
		}

		return nil
	}

	if err = afero.Walk(srcFs, _name, walker); err != nil {
		return 0, errs, err
	}

	return cnt, errs, nil
}

func (s *Controller) DirReaddir(name string) ([]map[string]interface{}, error) {

	_name := formatPath(name)

	afs := afero.Afero{s.fs}

	st, err := afs.ReadDir(_name)
	if err != nil {
		return nil, err
	}

	// preload data with correct size
	data := make([]map[string]interface{}, len(st))

	var ignore bool
	cnt := 0

	for i := range st {

		ignore = false

		// match against ignored files
		for _, match := range files_ignore {
			if strings.HasSuffix(st[i].Name(), match) {
				ignore = true
				break
			}
		}

		if ignore {
			continue
		}

		data[cnt] = map[string]interface{}{
			"size": st[i].Size(),
			"time": st[i].ModTime().UnixNano(),
			"dir":  st[i].IsDir(),
		}

		cnt++
	}

	// return data with real length
	return data[:cnt], err
}

func (s *Controller) DirCopy(name string, name2 string) error {

	_name := formatPath(name)
	_name2 := formatPath(name2)

	_, _, err := s.copier(s.fs, s.fs, _name, _name2)

	return err
}

func (s *Controller) DirDuplicate(name string, name2 string) error {

	return s.DirCopy(name, name2)
}

func (s *Controller) DirRename(name string, name2 string) error {

	_name := formatPath(name)
	_name2 := formatPath(name2)

	return s.fs.Rename(_name, _name2)
}

func (s *Controller) DirMove(name string, name2 string) error {

	return s.DirRename(name, name2)
}

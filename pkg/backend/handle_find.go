package backend

import (
	"os"
	"regexp"

	uerr "dynexo.de/ufyle/pkg/err"
	"dynexo.de/ufyle/pkg/io"

	"github.com/gobwas/glob"
	"github.com/spf13/afero"
)

type (
	Matcher interface {
		MatchString(string) bool
	}
)

func (s *Controller) findGlob(path string, matcher Matcher, limit int) ([]io.FileInfo, error) {

	var err error
	var st []io.FileInfo

	_path := formatPath(path)
	with_match := matcher != nil

	if limit <= 0 {
		limit = max_files
	}

	st = make([]io.FileInfo, limit)
	i := 0

	walker := func(path string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if with_match && matcher.MatchString(path) {

			st[i] = io.FileInfo{
				//name:    fi.Name(),
				FileName:    path,
				FileSize:    fi.Size(),
				FileModtime: fi.ModTime(),
				FileMode:    fi.Mode(),
			}

			i++
		}

		if i >= limit {
			//i--
			return uerr.ErrSizeLimit
		}

		return nil
	}

	if err = afero.Walk(s.fs, _path, walker); err != nil && err != uerr.ErrSizeLimit {
		return nil, err
	}

	return st[:i], nil
}

type WrapMatcher struct {
	matcher glob.Glob
}

func (s *WrapMatcher) MatchString(pattern string) bool {
	return s.matcher.Match(pattern)
}

func (s *Controller) FindGlob(path, pattern string, limit int) ([]string, error) {

	matcher := glob.MustCompile(pattern)

	list, err := s.findGlob(path, &WrapMatcher{matcher}, limit)
	list_len := len(list)

	entries := make([]string, len(list))

	for i := 0; i < list_len; i++ {
		entries[i] = list[i].String()
	}

	return entries, err
}

func (s *Controller) FindRegex(path, pattern string, limit int) ([]string, error) {

	matcher := regexp.MustCompile(pattern)

	list, err := s.findGlob(path, matcher, limit)
	list_len := len(list)

	entries := make([]string, len(list))

	for i := 0; i < list_len; i++ {
		entries[i] = list[i].String()
	}

	return entries, err
}

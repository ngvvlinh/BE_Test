package service

import (
	"testing"

	"dynexo.de/pkg/log"
	"dynexo.de/pkg/vfs"
	"github.com/spf13/afero"
)

func load() *Controller {

	s := NewController(log.New())

	osFs := afero.NewOsFs()

	if err := osFs.MkdirAll("vfs/fs1/", 0777); err != nil {
		return nil
	}

	if err := osFs.MkdirAll("vfs/fs2/", 0777); err != nil {
		return nil
	}

	fs, err := vfs.NewPrefixFs(afero.NewOsFs(), "vfs/fs1/")
	if err != nil {
		return nil
	}

	tempFs, err := vfs.NewPrefixFs(afero.NewOsFs(), "vfs/fs2/")
	if err != nil {
		return nil
	}

	s.LoadFs(fs, tempFs, "vfs/fs2/")

	return s
}

func aaTestReadddir(t *testing.T) {

	s := load()

	if s == nil {
		t.Error("Failed to load controller")
		t.FailNow()
	}

	fi, err := s.readddir(".", -1)
	if err != nil {
		t.FailNow()
	}

	for i := range fi {
		t.Logf("%s, %d, %v", fi[i].Name(), fi[i].Size(), fi[i].IsDir())
	}
}

func TestCopier1(t *testing.T) {

	s := load()

	if s == nil {
		t.Error("Failed to load controller")
		t.FailNow()
	}

	_, errs, err := s.copier(s.fs, s.tempFs, ".", ".")
	if err != nil {
		t.FailNow()
	}

	for i := range errs {
		t.Logf("%s, %s", errs[i].Op, errs[i].Path)
	}
}

func TestFindGlob(t *testing.T) {

	s := load()

	if s == nil {
		t.Error("Failed to load controller")
		t.FailNow()
	}

	list, err := s.FindGlob(".", "*marken*", 0)
	if err != nil {
		t.FailNow()
	}

	for i := range list {
		t.Logf("%v, %s", list[i].IsDir(), list[i].Path())
	}
}

func TestFindRegex(t *testing.T) {

	s := load()

	if s == nil {
		t.Error("Failed to load controller")
		t.FailNow()
	}

	list, err := s.FindRegex(".", "(?:^Briefm|.pdf$)", 4)
	if err != nil {
		t.FailNow()
	}

	for i := range list {
		t.Logf("%v, %s", list[i].IsDir(), list[i].Path())
	}
}

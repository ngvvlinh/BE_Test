package service

import (
	"sync"

	"dynexo.de/pkg/adt"
	"dynexo.de/pkg/io"

	"dynexo.de/pkg/log"
	"dynexo.de/tools"

	"github.com/kr/pretty"
	"github.com/spf13/afero"
)

var (
	dbg = pretty.Println
)

type (
	Controller struct {
		name    string
		closeCh *tools.CloseCh
		logger  log.ILogger

		fs      afero.Fs
		tempFs  afero.Fs
		basedir string

		bpool *io.BufferPool

		messages   *adt.PriorityQueue
		messagesMu sync.RWMutex

		mu sync.RWMutex
	}
)

func NewController(logger log.ILogger) *Controller {
	s := &Controller{}

	s.closeCh = tools.NewCloseCh("Controller", s.onClose)
	s.logger = logger
	s.bpool = io.NewBufferPool(32, 32, 0)

	return s
}

func (s *Controller) Debug() {

}

func (s *Controller) Log() log.ILogger {
	return s.logger
}

func (s *Controller) LoadFs(fs, tempFs afero.Fs, dir string) {
	s.fs = fs
	s.tempFs = tempFs
	s.basedir = dir
}

func (s *Controller) Fs() afero.Fs {
	return s.fs
}

func (s *Controller) FsRead() afero.Fs {
	return s.fs
}

func (s *Controller) FsWrite() afero.Fs {
	return s.fs
}

func (s *Controller) onClose() {

}

func (s *Controller) serve() {

	// serve MUST NO be blocking!!!

}

func (s *Controller) Serve() {
	s.serve()
}

func (s *Controller) ServeAndWait(cb func()) {

	s.serve()

	tools.WaitForCtrlC(cb)
}

func (s *Controller) AddMessage(str string) {

	s.messagesMu.RLock()
	defer s.messagesMu.RUnlock()

	s.messages.Push(str)
}

func (s *Controller) Messages() []string {

	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	if s.messages.Len() < 1 {
		return []string{}
	}

	str := make([]string, s.messages.Len())

	// FIXME we return string array without cast checking, but this should be safe

	for i := 0; i < s.messages.Len(); i++ {
		str[i] = s.messages.Pop().(string)
	}

	return str
}

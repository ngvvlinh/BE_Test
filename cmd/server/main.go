package main

import (
	"dynexo.de/pkg/vfs"
	"dynexo.de/ufyle/pkg/backend"
)

func main() {

	logger := getLogger()
	defer logger.Sync()

	// create new service controller
	s := backend.NewController(logger)
	s.LoadFs(vfs.OsFs(), "")

	backend.ListenAndServeService(":8443", s, logger)
}

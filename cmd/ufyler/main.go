package main

import (
	"dynexo.de/pkg/vfs"

	"dynexo.de/ufyle/internal/api"
	"dynexo.de/ufyle/pkg/service"

	"github.com/spf13/afero"
)

func main() {

	loadArgs()

	logger := getLogger()
	defer logger.Sync()

	// create new service controller
	s := service.NewController(logger)

	s.Serve()

	if err := loadFs(s); err != nil {
		logger.Panicf("Failed to load filesystem", err)
	}

	// create new api service
	c := api.NewServiceApi(s, logger)

	_addr := argConfigAddr

	server, err := api.NewApiServer(_addr, logger)
	if err != nil {
		logger.Panicf("Failed to start service listener on %s", _addr)
	}

	// run the api service
	c.Run(server)
}

func loadFs(s *service.Controller) error {

	osFs := afero.NewOsFs()

	if err := osFs.MkdirAll(baseDir, 0777); err != nil {
		return err
	}

	if err := osFs.MkdirAll(tempDir, 0777); err != nil {
		return err
	}

	fs, err := vfs.NewPrefixFs(afero.NewOsFs(), baseDir)
	if err != nil {
		return err
	}

	tempFs, err := vfs.NewPrefixFs(afero.NewOsFs(), tempDir)
	if err != nil {
		return err
	}

	s.LoadFs(fs, tempFs, tempDir)

	return nil
}

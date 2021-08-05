// +build !debug

package main

import (
	"dynexo.de/pkg/log"
)

func getLogger() log.ILogger {
	return log.New("ufyler")
}

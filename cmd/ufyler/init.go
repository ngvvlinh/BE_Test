package main

import (
	"flag"
	"fmt"

	_ "net/http/pprof"

	"dynexo.de/product"
)

const (
	baseDir = "d:/temp/ufyle/vfs"
	tempDir = "d:/temp/ufyle/tempFs"
)

var (
	argConfigFile string
	argConfigAddr string
	argIsDebug    bool
	argIsPProf    bool
)

func init() {
	flag.StringVar(&argConfigFile, "config", "config.json", "Configuration file")
	flag.StringVar(&argConfigAddr, "addr", ":889", "Listener address")
	flag.BoolVar(&argIsDebug, "debug", false, "Enable debug output")
	flag.BoolVar(&argIsPProf, "pprof", false, "Enable profiling")

}

func loadArgs() {
	flag.Parse()
}

func update() {
	if done, _ := product.SelfUpdate(); done {
		fmt.Println("Software version updated. Please start again.")
		return
	}
}

func loadPprof() {

}

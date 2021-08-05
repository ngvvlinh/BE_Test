// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", ".")
	return cmd.Run()
}

func Dev() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")

	// "-o", APP_NAME,

	cmd := exec.Command("go", "build", "-v", "-race", "-tags", "debug", ".")
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	//return os.Rename("./MyApp", "/usr/bin/MyApp")
	return nil
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	//	fmt.Println("Installing Deps...")
	//	cmd := exec.Command("go", "get", "github.com/stretchr/piglatin")
	//	return cmd.Run()
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("ufyler.exe")
	os.RemoveAll("ufyler")
	os.RemoveAll("ufyler-dev.exe")
	os.RemoveAll("ufyler-dev")
}

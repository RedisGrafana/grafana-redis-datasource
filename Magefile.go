//+build mage

package main

import (
	"os"
	"path/filepath"

	// mage:import
	build "github.com/grafana/grafana-plugin-sdk-go/build"
	"github.com/magefile/mage/sh"
)

// runs backend tests and makes a txt coverage report in "atomic" mode and html coverage report.
func Cover() error {
	// Create a coverage file if it does not already exist
	if err := os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm); err != nil {
		return err
	}

	if err := sh.RunV("go", "test", "./pkg/...", "-v", "-cover", "-covermode=atomic", "-coverprofile=coverage/backend.txt"); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.txt", "-o", "coverage/backend.html"); err != nil {
		return err
	}

	return nil
}

// Default configures the default target.
var Default = build.BuildAll

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

// up docker-compose environment and run backend tests including integration tests with coverage report
func Integration() error {
	// Create a coverage file if it does not already exist
	if err := os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm); err != nil {
		return err
	}

	if err := Up(); err != nil {
		return err
	}

	defer Down()

	if err := sh.RunV("go", "test", "./pkg/...", "-tags=integration", "-v", "-cover", "-covermode=atomic", "-coverprofile=coverage/backend.txt"); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.txt", "-o", "coverage/backend.html"); err != nil {
		return err
	}

	return nil
}

// up docker-compose environment from integration tests
func Up() error {
	return sh.RunV("docker-compose", "-f", "docker-compose-test.yml", "-p", "grd-integration", "up", "-d")
}

// down docker-compose environment from integration tests
func Down() error {
	return sh.RunV("docker-compose", "-f", "docker-compose-test.yml", "-p", "grd-integration", "down")
}

// Default configures the default target.
var Default = build.BuildAll

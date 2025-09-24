//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
)

// Default target to run when none is specified
// in preference to running all.
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(BuildGos, BuildGosc)
	return nil
}

// A custom install step that runs after build
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	err := os.Rename("./gos", os.ExpandEnv("$HOME/go/bin/gos"))
	if err != nil {
		return err
	}
	return os.Rename("./gosc", os.ExpandEnv("$HOME/go/bin/gosc"))
}

// A step on it's own
func BuildGos() error {
	fmt.Println("Building gos...")
	cmd := exec.Command("go", "build", "-o", "gos", "cmd/gos/main.go")
	return cmd.Run()
}

// A step on it's own
func BuildGosc() error {
	fmt.Println("Building gosc...")
	cmd := exec.Command("go", "build", "-o", "gosc", "cmd/gosc/main.go")
	return cmd.Run()
}

// Run gos
func Run() error {
	mg.Deps(Dev)
	fmt.Println("Running...")
	cmd := exec.Command("go", "run", "cmd/gos/main.go")
	return cmd.Run()
}

// Run development steps
func Dev() error {
	mg.Deps(Test, Vet, Lint)
	fmt.Println("Building dev...")
	cmd := exec.Command("go", "build", "-race", "-o", "gos", "cmd/gos/main.go")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "build", "-race", "-o", "gosc", "cmd/gosc/main.go")
	return cmd.Run()
}

// Run tests
func Test() error {
	fmt.Println("Testing...")
	cmd := exec.Command("go", "clean", "-testcache")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "test", "-v", "./...")
	return cmd.Run()
}

// Run fuzz tests
func Fuzz() error {
	fmt.Println("Fuzzing...")
	cmd := exec.Command("go", "clean", "-testcache")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "test", "./internal/entry/", "-fuzz=FuzzExtractURLs", "-fuzztime=10s")
	return cmd.Run()
}

// Run vet
func Vet() error {
	fmt.Println("Vetting...")
	cmd := exec.Command("go", "vet", "./...")
	return cmd.Run()
}

// Run lint
func Lint() error {
	fmt.Println("Linting...")
	cmd := exec.Command("golangci-lint", "run")
	return cmd.Run()
}

// Install dev tools
func DevInstall() error {
	fmt.Println("Installing dev tools...")
	cmd := exec.Command("go", "install", "golang.org/x/tools/gopls@latest")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	return cmd.Run()
}

// Clean up after build
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("gos")
	os.RemoveAll("gosc")
}

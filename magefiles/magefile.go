package main

import (
	"fmt"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var appExecutable = "bin/go-nats-app"
var appDir = "."

const goCompiler = "go"
const gopherjsCompiler = "gopherjs"

var appGlobs = []string{
	"magefiles/magefile.go",
	"*.go",
	"*/*.go",
}

//goland:noinspection GoBoolExpressions
func init() {
	if runtime.GOOS == "windows" {
		appExecutable += ".exe"
	}
}

func buildApp() error {
	changes, err := target.Glob(appExecutable, appGlobs...)
	if err != nil {
		return err
	}
	if !changes {
		return nil
	}
	fmt.Println("> Building App...")
	return sh.RunV(goCompiler, "build", "-o", appExecutable, appDir)
}

func BuildFrontend() error {
	changes, err := target.Glob("web/app.js", appGlobs...)
	if err != nil {
		return err
	}
	if !changes {
		return nil
	}
	fmt.Println("> Building Frontend...")
	return sh.RunV(gopherjsCompiler, "build", "-m", "-o", "web/app.js", appDir)
}

func Build() error {
	mg.Deps(BuildFrontend)
	return buildApp()
}

func Run() error {
	mg.Deps(Build)
	return sh.RunV(appExecutable)
}

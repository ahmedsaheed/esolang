package main

import (
	_ "embed"
	"esolang/lang-esolang/repl"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

//go:embed stdlib/array.eso
var Arraylib string

//go:embed stdlib/bool.eso
var BoolLib string
var lib = Arraylib + "\n" + BoolLib

// build using: go build -tags netgo -ldflags '-s -w' -o esolang
func main() {
	replMode := flag.Bool("repl", false, "Start the repl")
	logger := log.New(os.Stderr)
	flag.Parse()

	if *replMode {
		repl.Start(os.Stdin, os.Stdout, lib)
	}

	if len(flag.Args()) > 0 {
		file := flag.Args()[0]
		inputFile, err := os.ReadFile(file)

		if err != nil {
			logger.Error("Error reading file %s", os.Args[1])
			os.Exit(1)
		}
		fileExtension := filepath.Ext(file)
		if fileExtension != ".eso" {
			logger.Error("Invalid file extension. Please provide a file with .eso extension")
			os.Exit(1)
		}

		repl.Execute(string(inputFile), lib)
	} else {
		logger.Warn("No file provided. Please provide a file to run or use the -repl flag to start the repl.")
		logger.Warn("Usage: esolang <path-to-filename>")
		logger.Warn("Usage: esolang -repl")
		os.Exit(1)
	}

}

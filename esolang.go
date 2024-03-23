package main

import (
	"esolang/lang-esolang/repl"
	"flag"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func main() {
	replMode := flag.Bool("repl", false, "Start the repl")
	logger := log.New(os.Stderr)
	flag.Parse()

	if *replMode {
		repl.Start(os.Stdin, os.Stdout)
	}

	if len(flag.Args()) > 0 {
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

		repl.Execute(string(inputFile))
	} else {
		logger.Warn("No file provided. Please provide a file to run or use the -repl flag to start the repl.")
		logger.Warn("Usage: esolang <path-to-filename>")
		logger.Warn("Usage: esolang -repl")
		os.Exit(1)
	}

}

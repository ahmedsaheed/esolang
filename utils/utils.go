package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var SearchPaths []string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting cwd: %s", err)
	}

	if e := os.Getenv("ESOPATH"); e != "" {
		tokens := strings.Split(e, ":")
		for _, token := range tokens {
			AddPath(token) // ignore errors
		}
	} else {
		SearchPaths = append(SearchPaths, cwd)
	}
}

func AddPath(path string) error {
	path = os.ExpandEnv(filepath.Clean(path))
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	SearchPaths = append(SearchPaths, absPath)
	return nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsBuiltinModule(name string) bool {
	return map[bool]bool{true: true, false: false}[strings.HasPrefix(name, "eso/") && strings.Count(name, "/") == 1]
}

func FindModule(name string) string {
	basename := fmt.Sprintf("%s.eso", name)
	for _, p := range SearchPaths {
		filename := filepath.Join(p, basename)
		if Exists(filename) {
			return filename
		}
	}
	return ""
}

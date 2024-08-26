package main

// FIXME:
// - Fix the normal ls with no flags

// Allowed imports:
// "fmt"
// "os"
// "os/user"
// "strconv"
// "strings"
// "syscall"
// "time"
// "math/rand"
// "errors"
// "io/fs"

import (
	"fmt"
	"os"

	RUN "ok/my_ls/PGK"
)

// Define command-line flags
var (
	longListing bool
	recursive   bool
	allFiles    bool
	reverse     bool
	sortByTime  bool
	MyLine      []string
)

// ANSI escape codes for colors
const (
	Blue  = "\033[34m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

func main() {
	// Parse command-line flags
	MyLine, longListing, recursive, allFiles, reverse, sortByTime = RUN.ParseFlags(os.Args[1:])

	// Determine paths to list files from
	var paths []string
	if len(MyLine) == 0 {
		paths = []string{"./"} // Default to current directory
	} else {
		paths = MyLine
	}

	paths = RUN.CheckPathsExist(paths)

	NewFileSlice, paths := RUN.TestFile(paths)

	if NewFileSlice != nil {
		RUN.PrintFiles(NewFileSlice, longListing, reverse, sortByTime)

	}

	// Iterate through each specified path
	for _, path := range paths {
		RUN.HandlePath(paths, path, longListing, recursive, allFiles, reverse, sortByTime)
		fmt.Println()
	}
}

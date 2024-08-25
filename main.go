package main

// TODO:
// - Adjust the print in the long listing and other places if needed (Not sure if it is important)

// FIXME:
// - Fix the -l flag (1- The adjesment)

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

	// Iterate through each specified path
	for _, path := range paths {
		RUN.HandlePath(path, longListing, recursive, allFiles, reverse, sortByTime)
	}
}

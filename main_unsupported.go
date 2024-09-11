//go:build !linux

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "dirscope: Linux is required; on Windows, run DirScope inside WSL Ubuntu")
	os.Exit(2)
}

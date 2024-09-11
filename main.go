//go:build linux

package main

import (
	"os"

	"github.com/ihamzaihsan/dirscope/internal/listing"
)

func main() {
	os.Exit(listing.Run(os.Args[1:], os.Stdout, os.Stderr))
}

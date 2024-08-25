package main

import (
	"fmt"
	"os"
	"syscall"
)

// CalculateTotalBlocks calculates the total number of disk blocks used by files in the specified directory.
func CalculateTotalBlocks(path string) (int64, error) {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer dir.Close()

	// Read directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return 0, err
	}

	var totalBlocks int64

	// Iterate through the files and calculate the total number of blocks
	for _, fileInfo := range fileInfos {
		if stat, ok := fileInfo.Sys().(*syscall.Stat_t); ok {
			totalBlocks += stat.Blocks
		}
	}

	return totalBlocks, nil
}

func main() {
	// Specify the directory path
	path := "./"

	// Calculate total blocks
	totalBlocks, err := CalculateTotalBlocks(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// The total number of blocks used by the files
	fmt.Printf("total %d\n", totalBlocks)
}

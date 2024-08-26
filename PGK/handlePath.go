package RUN

import (
	"fmt"
	"os"
)

// HandlePath processes the files and directories at the given path
func HandlePath(paths []string, path string, longListing, recursive, allFiles, reverse, sortByTime bool) {

	// Continue only if the path exists
	if len(paths) == 0 {
		fmt.Println("Error: No valid paths provided")
		return
	}

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer dir.Close()

	// Read directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Sort files by name
	SortByName(fileInfos)

	// Handle hidden files based on the -a flag
	if allFiles {
		// Add current and parent directory (".", "..") if -a is specified
		currentDir, _ := os.Stat(".")
		parentDir, _ := os.Stat("..")
		fileInfos = append([]os.FileInfo{currentDir, parentDir}, fileInfos...)
	} else {
		fileInfos = FilterHiddenFiles(fileInfos)
	}

	// Sort files by modification time if the -t flag is specified
	if sortByTime {
		SortByModTime(fileInfos)
	}

	// Reverse the order of files if the -r flag is specified
	if reverse {
		HandleReverse(fileInfos)
	}

	// Print files in long listing format if the -l flag is specified
	if longListing && !recursive {
		HandleLongListing(fileInfos, true)
	} else if recursive {
		// Recursively list subdirectories if the -R flag is specified
		HandleRecursive(path, longListing, allFiles, reverse, sortByTime)
	} else {
		// Print file names if not long listing or recursive
		if len(paths) > 1 {
			for _, file := range fileInfos {
				if file.IsDir() {
					// PrintFileName(file)

					fmt.Println(file.Name(), ":")
				}
				PrintFileName(file)
			}
			fmt.Println()
		} else {
			for _, file := range fileInfos {
				PrintFileName(file)
			}
		}
	}

	// fmt.Println() // Print a new line after listing files
}

// CheckPathsExist checks if each path in the list exists and removes non-existing ones.
func CheckPathsExist(paths []string) []string {
	var validPaths []string

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("ls: cannot access '%s': No such file or directory\n", path)
		} else {
			validPaths = append(validPaths, path)
		}
	}

	return validPaths
}


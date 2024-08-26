package RUN

import (
	"fmt"
	"os"
)

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		// If there's an error, return false
		return false
	}
	return info.Mode().IsRegular()
}

func TestFile(paths []string) ([]string, []string) {
	var NewFileSlice []string
	for _, path := range paths {
		if IsFile(path) {
			NewFileSlice = append(NewFileSlice, path)
			paths = RemoveFromSlice(paths, path)
		}
	}
	return NewFileSlice, paths

}

func PrintFiles(NewFileSlice []string, longListing, reverse, sortByTime bool) {
	var T []os.FileInfo
	for _, file := range NewFileSlice {
		// Get the FileInfo for the provided path
		info, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		T = append(T, info)
	}
	SortByName(T)

	if sortByTime {
		SortByModTime(T)
	}

	if reverse {
		HandleReverse(T)
	}

	if longListing {
		HandleLongListing(T, false)
	}
	for _, ok := range T {
		PrintFileName(ok)
	}
	fmt.Println()
	fmt.Println()

}

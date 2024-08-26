//go:build windows
// +build windows

package RUN

import (
	"fmt"
	"os"
)

func RunPlatformSpecificCode() {
	fmt.Println("This is Windows-specific code.")
	// Add Windows-specific code here
}

func CalculateTotalSize(path string) (int64, error) {
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

	var totalSize int64

	// Iterate through the files and calculate the total size
	for _, fileInfo := range fileInfos {
		totalSize += fileInfo.Size()
	}

	return totalSize, nil
}

// Not importatn to fix
func GetGroupName(path string) (string, error) {
	// // Get the file information
	// fileInfo, err := os.Stat(path)
	// if err != nil {
	// 	return "", err
	// }

	// // Extract the underlying data interface (*syscall.Stat_t)
	// stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	// if !ok {
	// 	return "", fmt.Errorf("unable to retrieve file info")
	// }

	// // Lookup the group name using the GID
	// group, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	// if err != nil {
	// 	return "", err
	// }

	// return group.Name, nil

	return "", nil
}

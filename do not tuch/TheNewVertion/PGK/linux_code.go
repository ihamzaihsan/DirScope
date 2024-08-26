//go:build linux
// +build linux

package RUN

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
)

func userName() string {
	currentUser, _ := user.Current()
	fullUser := currentUser.Username
	last := strings.Split(fullUser, "\\")
	return last[len(last)-1]
}

func CalculateMaxWidths(files []os.FileInfo) (nlinkWidth, userWidth, groupWidth, sizeWidth int) {
	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t)

		// Calculate widths for number of links, username, group name, and file size
		if nlinkLen := len(fmt.Sprintf("%d", stat.Nlink)); nlinkLen > nlinkWidth {
			nlinkWidth = nlinkLen
		}
		if userLen := len(userName()); userLen > userWidth {
			userWidth = userLen
		}
		if groupLen := len(groupName(file)); groupLen > groupWidth {
			groupWidth = groupLen
		}
		if sizeLen := len(fmt.Sprintf("%d", file.Size())); sizeLen > sizeWidth {
			sizeWidth = sizeLen
		}
	}
	return nlinkWidth, userWidth, groupWidth, sizeWidth
}

func RunPlatformSpecificCode() {
	fmt.Println("This is Linux-specific code.")
	// Add Linux-specific code here
}

func GetHardLinkCount(path string) (uint64, error) {
	// Get file info
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return 0, err
	}

	// Convert the file info to a syscall.Stat_t struct to access the number of links
	stat := fileInfo.Sys().(*syscall.Stat_t)

	// Return the number of hard links
	return stat.Nlink, nil
}

func CalculateTotalBlocks(files []os.FileInfo) (int64, error) {

	var totalBlocks int64

	// Iterate through the files and calculate the total number of blocks
	for _, fileInfo := range files {
		
		stat, _ := fileInfo.Sys().(*syscall.Stat_t)
		
		totalBlocks = totalBlocks + stat.Blocks


	}
	totalBlocks = (totalBlocks / 2) 

	return totalBlocks, nil
}

func GetGroupName(path string) (string, error) {
	// Get the file information
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// Extract the underlying data interface (*syscall.Stat_t)
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("unable to retrieve file info")
	}

	// Lookup the group name using the GID
	group, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		return "", err
	}

	return group.Name, nil
}

func groupName(fileInfo os.FileInfo) string {
		// Extract the underlying data interface (*syscall.Stat_t)
		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return ""
		}
	
		// Lookup the group name using the GID
		group, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
		if err != nil {
			return ""
		}
	
		return group.Name
}

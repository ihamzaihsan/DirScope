package RUN

import (
	"fmt"
	"os"
	// "os/user"
	"strings"
	"syscall"
)

var (
	longListing bool
	reverse bool
	MyLine []string
)

// ANSI escape codes for colors
const (
	Blue  = "\033[34m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

// SortByModTime sorts files and folders by modification time, with the newest first.
func SortByModTime(fileInfos []os.FileInfo) {
	n := len(fileInfos)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if fileInfos[j].ModTime().Before(fileInfos[j+1].ModTime()) {
				fileInfos[j], fileInfos[j+1] = fileInfos[j+1], fileInfos[j]
			}
		}
	}
}

// SortByName sorts files and folders alphabetically by their names.
func SortByName(fileInfos []os.FileInfo) {
	n := len(fileInfos)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if strings.ToLower(fileInfos[j].Name()) > strings.ToLower(fileInfos[j+1].Name()) {
				fileInfos[j], fileInfos[j+1] = fileInfos[j+1], fileInfos[j]
			}
		}
	}
}

// handleLongListing prints files in a long listing format
func HandleLongListing(files []os.FileInfo, path string) {
	totalBlocks, err := CalculateTotalBlocks(files)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// var num int64

	// for _, f := range files {
	// 	stat := f.Sys().(*syscall.Stat_t)
	// 	num = num + stat.Blocks
	// }

	fmt.Println("total", totalBlocks)

	// currentUser, err := user.Current()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// GroupNeam, _ := GetGroupName(path)

	// fullUser := currentUser.Username
	// Last := strings.Split(fullUser, "\\")
	// USER := Last[len(Last)-1]

	// Calculate max widths for columns
	nlinkWidth, userWidth, groupWidth, sizeWidth := CalculateMaxWidths(files)

	//FIXME: the ajecment
	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t)

		// fmt.Printf("%s %d %- %s %s %d %s ",
		// 	file.Mode().String(),
		// 	stat.Nlink,
		// 	USER,
		// 	GroupNeam,
		// 	file.Size(),
		// 	file.ModTime().Format("Jan 2 15:04"),
		// )
		// PrintFileName(file)
		// fmt.Println()

		fmt.Printf("%s %*d %-*s %-*s %*d %s ",
			file.Mode().String(),
			nlinkWidth, stat.Nlink,
			userWidth, userName(),
			groupWidth, groupName(file),
			sizeWidth, file.Size(),
			file.ModTime().Format("Jan 2 15:04"),
		)
		PrintFileName(file)
		fmt.Println()

	}

// 	fmt.Printf("%s %*d %-*s %-*s %*d %s ",
// 	file.Mode().String(),
// 	nlinkWidth, stat.Nlink,
// 	userWidth, userName(),
// 	groupWidth, groupName(),
// 	sizeWidth, file.Size(),
// 	file.ModTime().Format("Jan 2 15:04"),
// )
// PrintFileName(file)
// fmt.Println()
}

func HandleRecursive(path string, longListing, allFiles, reverse, sortByTime bool) {
	// Get the FileInfo for the provided path
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check if the path is a directory
	if info.IsDir() {
		// Open the directory
		dir, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer dir.Close()

		// Read the contents of the directory
		files, err := dir.Readdir(-1)
		if err != nil {
			fmt.Println(err)
			return
		}

		var dirFileInfos []os.FileInfo

		// Handle the -a flag to include hidden files, . and ..
		if allFiles {
			// Manually add "." and ".."
			currentDir, _ := os.Stat(path + "/.")
			parentDir, _ := os.Stat(path + "/..")
			dirFileInfos = append([]os.FileInfo{currentDir, parentDir}, files...)
		} else {
			dirFileInfos = FilterHiddenFiles(files)
		}

		// Sort files by name or time
		SortByName(dirFileInfos)
		if sortByTime {
			SortByModTime(dirFileInfos)
		}

		// Reverse the order if the flag is set
		if reverse {
			HandleReverse(dirFileInfos)
		}

		// Print the files in long listing format if the flag is set
		if longListing {
			HandleLongListing(dirFileInfos, path)
		} else {
			for _, file := range dirFileInfos {
				PrintFileName(file)
			}
			fmt.Println()
		}

		// Recursively handle subdirectories
		for _, file := range dirFileInfos {
			if file.IsDir() && file.Name() != "." && file.Name() != ".." {
				fmt.Printf("\n%s:\n", path+"/"+file.Name())
				HandleRecursive(path+"/"+file.Name(), longListing, allFiles, reverse, sortByTime)
			}
		}
	} else {
		PrintFileName(info)
	}
}

// FilterHiddenFiles filters out hidden files
func FilterHiddenFiles(files []os.FileInfo) []os.FileInfo {
	var filteredFiles []os.FileInfo
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

// HandleReverse reverses the order of files
func HandleReverse(files []os.FileInfo) {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
}

// printFileName prints the file name, coloring directories blue, .exe files green, and adding a backslash or asterisk at the end
func PrintFileName(file os.FileInfo) {
	if file.IsDir() {
		fmt.Printf("%s%s%s/%s  ", Blue, file.Name(), Reset, Reset)
	} else if strings.HasSuffix(file.Name(), ".exe") {
		fmt.Printf("%s%s%s*%s  ", Green, file.Name(), Reset, Reset)
	} else {
		fmt.Print(file.Name(), "  ")
	}
}

package RUN

import (
	"fmt"
	"os"
	"strings"
)

func ParseFlags(MyLine []string) ([]string, bool, bool, bool, bool, bool) {
	var longListing, recursive, allFiles, reverse, sortByTime bool
	for _, arg := range MyLine {
		if arg == "--l" {
			MyLine = RemoveFromSlice(MyLine, "--l")
			continue
		}
		switch arg {
		case "-l":
			longListing = true
			MyLine = RemoveFromSlice(MyLine, "-l")
		case "-R":
			recursive = true
			MyLine = RemoveFromSlice(MyLine, "-R")
		case "-a":
			allFiles = true
			MyLine = RemoveFromSlice(MyLine, "-a")
		case "-r":
			reverse = true
			MyLine = RemoveFromSlice(MyLine, "-r")
		case "-t":
			sortByTime = true
			MyLine = RemoveFromSlice(MyLine, "-t")
		default:
			 if strings.HasPrefix(arg, "--") {
				fmt.Printf("my_ls: invalid option '%s'\n", arg)
				fmt.Println("Try 'my_ls --help' for more information.")
				os.Exit(1)

			} else if strings.HasPrefix(arg, "-") {

				for i, ok := range arg {
					if i == 0 {
					} else if ok == 'l' {
						longListing = true

					} else if ok == 'R' {
						recursive = true

					} else if ok == 'a' {
						allFiles = true

					} else if ok == 'r' {
						reverse = true

					} else if ok == 't' {
						sortByTime = true
						
					} else {
						fmt.Printf("my_ls: invalid option '%s'\n", string(ok))
						fmt.Println("Try 'my_ls --help' for more information.")
						os.Exit(1)
					}
				}
				MyLine = RemoveFromSlice(MyLine, arg)

			}

		}
	}
	// for _, arg := range MyLine {
		// if strings.HasPrefix(arg, "-") {
	// 		fmt.Printf("my_ls: invalid option '%s'\n", arg)
	// 		fmt.Println("Try 'my_ls --help' for more information.")
	// 		os.Exit(1)
	// 	}
	// }

	// if strings.HasPrefix(arg, "---") {
	// 	fmt.Printf("my_ls: invalid option '%s'\n", arg)
	// 	fmt.Println("Try 'my_ls --help' for more information.")
	// 	os.Exit(1)
	// }
	return MyLine, longListing, recursive, allFiles, reverse, sortByTime
}

func RemoveFromSlice(slice []string, s string) []string {
	var result []string
	for _, v := range slice {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}

func replaceDoubleDash(statement string) string {
    if strings.HasPrefix(statement, "--") {
        return "-" + statement[2:] // Replace "--" with "-"
    }
    return statement
}

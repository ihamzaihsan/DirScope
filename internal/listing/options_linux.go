package listing

import (
	"fmt"
	"strings"
)

const help = `DirScope: Command-Line File Listing Tool
Usage: dirscope [options] [path ...]

List the current directory when no paths are supplied.
Options can be combined (for example, -laR) or placed between paths.

  -l               Show permissions, links, owner, group, size and time
  -R               List subdirectories recursively
  -a               Include dotfiles, . and ..
  -r               Reverse the selected sort order
  -t               Sort by modification time, newest first
  --color=WHEN     Color names: auto (default), always or never
  --help           Show this help
  --               Treat remaining arguments as paths

Linux only. Names are sorted bytewise, matching LC_ALL=C ls.
`

type options struct {
	long, recursive, all, reverse, byTime bool
	color                                 string
}

func parse(args []string) (options, []string, bool, error) {
	opts := options{color: "auto"}
	var paths []string
	endOptions, showHelp := false, false
	for _, arg := range args {
		if endOptions || arg == "-" || !strings.HasPrefix(arg, "-") {
			paths = append(paths, arg)
			continue
		}
		switch {
		case arg == "--":
			endOptions = true
		case arg == "--help":
			showHelp = true
		case strings.HasPrefix(arg, "--color="):
			opts.color = strings.TrimPrefix(arg, "--color=")
			if opts.color != "auto" && opts.color != "always" && opts.color != "never" {
				return opts, nil, false, fmt.Errorf("invalid color mode %q", opts.color)
			}
		case strings.HasPrefix(arg, "--"):
			return opts, nil, false, fmt.Errorf("unknown option %q", arg)
		default:
			for _, flag := range arg[1:] {
				switch flag {
				case 'l':
					opts.long = true
				case 'R':
					opts.recursive = true
				case 'a':
					opts.all = true
				case 'r':
					opts.reverse = true
				case 't':
					opts.byTime = true
				default:
					return opts, nil, false, fmt.Errorf("unknown option -%c", flag)
				}
			}
		}
	}
	if len(paths) == 0 {
		paths = []string{"."}
	}
	return opts, paths, showHelp, nil
}

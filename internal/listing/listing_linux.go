// Package listing implements DirScope's Linux directory listings.
package listing

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type entry struct {
	name, path string
	info       os.FileInfo
}

type runner struct {
	opts           options
	out, errOut    *os.File
	status         int
	printed, color bool
	now            time.Time
	timeStyle      string
	users, groups  map[uint32]string
}

// Run lists paths and returns 0 on success, 1 on directory-entry errors,
// or 2 on invalid arguments, inaccessible operands or output errors.
func Run(args []string, out, errOut *os.File) int {
	opts, paths, showHelp, err := parse(args)
	if err != nil {
		fmt.Fprintf(errOut, "dirscope: %v\nTry 'dirscope --help' for usage.\n", err)
		return 2
	}
	if showHelp {
		if _, err := fmt.Fprint(out, help); err != nil {
			return 2
		}
		return 0
	}
	r := runner{opts: opts, out: out, errOut: errOut, now: time.Now(),
		users: make(map[uint32]string), groups: make(map[uint32]string)}
	if opts.long {
		r.timeStyle = os.Getenv("TIME_STYLE")
		switch r.timeStyle {
		case "", "locale", "long-iso", "full-iso", "iso":
		default:
			fmt.Fprintf(errOut, "dirscope: unsupported TIME_STYLE %q; use locale, long-iso, full-iso or iso\n", r.timeStyle)
			return 2
		}
	}
	r.color = opts.color == "always"
	if opts.color == "auto" && os.Getenv("TERM") != "dumb" {
		if info, err := out.Stat(); err == nil {
			r.color = info.Mode()&os.ModeCharDevice != 0
		}
	}
	var files, dirs []entry
	for _, path := range paths {
		info, err := os.Lstat(path)
		// A directory symlink operand is followed in short mode.
		// A trailing slash forces resolution by the operating system.
		if err == nil && !opts.long && info.Mode()&os.ModeSymlink != 0 {
			if target, statErr := os.Stat(path); statErr == nil && target.IsDir() {
				info = target
			}
		}
		if err != nil {
			r.report(path, err, 2)
			continue
		}
		e := entry{name: path, path: path, info: info}
		if info.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	r.sort(files)
	r.sort(dirs)
	if len(files) > 0 {
		r.printEntries(files, false, dirs)
		r.printed = true
	}
	for _, dir := range dirs {
		r.directory(dir.path, opts.recursive || len(paths) > 1)
	}
	return r.status
}

func (r *runner) report(path string, err error, status int) {
	if status > r.status {
		r.status = status
	}
	fmt.Fprintf(r.errOut, "dirscope: %q: %v\n", path, err)
}

func (r *runner) write(format string, args ...any) {
	if _, err := fmt.Fprintf(r.out, format, args...); err != nil {
		if r.status != 2 {
			fmt.Fprintf(r.errOut, "dirscope: write output: %v\n", err)
		}
		r.status = 2
	}
}

// childPath preserves the user's spelling in recursive directory headings.
func childPath(parent, name string) string {
	if strings.HasSuffix(parent, "/") {
		return parent + name
	}
	return parent + "/" + name
}

func (r *runner) directory(path string, heading bool) {
	dir, err := os.Open(path)
	if err != nil {
		r.report(path, err, 2)
		return
	}
	names, readErr := dir.Readdirnames(-1)
	closeErr := dir.Close() // Close before descending to avoid retaining a descriptor per level.
	if readErr != nil {
		r.report(path, readErr, 2)
	}
	if closeErr != nil {
		r.report(path, closeErr, 1)
	}
	if r.opts.all {
		names = append(names, ".", "..")
	}
	var entries []entry
	for _, name := range names {
		if !r.opts.all && strings.HasPrefix(name, ".") {
			continue
		}
		fullPath := childPath(path, name)
		info, err := os.Lstat(fullPath)
		if err != nil {
			r.report(fullPath, err, 1)
			continue
		}
		entries = append(entries, entry{name: name, path: fullPath, info: info})
	}
	r.sort(entries)
	if heading {
		if r.printed {
			r.write("\n")
		}
		r.write("%s:\n", displayName(path))
	}
	r.printEntries(entries, true, nil)
	r.printed = true
	if r.opts.recursive {
		for _, e := range entries {
			// Lstat prevents traversal through symlinks and resulting cycles.
			if e.info.IsDir() && e.name != "." && e.name != ".." {
				r.directory(e.path, true)
			}
		}
	}
}

func (r *runner) before(a, b entry) bool {
	if r.opts.byTime && !a.info.ModTime().Equal(b.info.ModTime()) {
		if r.opts.reverse {
			return a.info.ModTime().Before(b.info.ModTime())
		}
		return a.info.ModTime().After(b.info.ModTime())
	}
	if r.opts.reverse {
		return a.name > b.name
	}
	return a.name < b.name
}

// Stable merge sort keeps sorting O(n log n) using the assignment's allowed imports.
func (r *runner) sort(entries []entry) {
	if len(entries) < 2 {
		return
	}
	buffer := make([]entry, len(entries))
	var merge func(int, int)
	merge = func(lo, hi int) {
		if hi-lo < 2 {
			return
		}
		mid := lo + (hi-lo)/2
		merge(lo, mid)
		merge(mid, hi)
		i, j := lo, mid
		for k := lo; k < hi; k++ {
			if i < mid && (j >= hi || !r.before(entries[j], entries[i])) {
				buffer[k] = entries[i]
				i++
			} else {
				buffer[k] = entries[j]
				j++
			}
		}
		copy(entries[lo:hi], buffer[lo:hi])
	}
	merge(0, len(entries))
}

package listing

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type row struct {
	entry                          entry
	mode, owner, group, size, date string
	links                          uint64
	blocks                         int64
	major, minor                   string
}

func (r *runner) owner(id uint32) string {
	if name, ok := r.users[id]; ok {
		return name
	}
	name := strconv.FormatUint(uint64(id), 10)
	if account, err := user.LookupId(name); err == nil {
		name = account.Username
	}
	r.users[id] = name
	return name
}

func (r *runner) group(id uint32) string {
	if name, ok := r.groups[id]; ok {
		return name
	}
	name := strconv.FormatUint(uint64(id), 10)
	if group, err := user.LookupGroupId(name); err == nil {
		name = group.Name
	}
	r.groups[id] = name
	return name
}

func permissions(mode uint32) string {
	text := []byte("----------")
	switch mode & syscall.S_IFMT {
	case syscall.S_IFDIR:
		text[0] = 'd'
	case syscall.S_IFLNK:
		text[0] = 'l'
	case syscall.S_IFBLK:
		text[0] = 'b'
	case syscall.S_IFCHR:
		text[0] = 'c'
	case syscall.S_IFIFO:
		text[0] = 'p'
	case syscall.S_IFSOCK:
		text[0] = 's'
	}
	for i, letter := range []byte("rwxrwxrwx") {
		if mode&(1<<uint(8-i)) != 0 {
			text[i+1] = letter
		}
	}
	for _, special := range []struct {
		bit     uint32
		index   int
		on, off byte
	}{
		{syscall.S_ISUID, 3, 's', 'S'}, {syscall.S_ISGID, 6, 's', 'S'}, {syscall.S_ISVTX, 9, 't', 'T'},
	} {
		if mode&special.bit != 0 {
			if text[special.index] == 'x' {
				text[special.index] = special.on
			} else {
				text[special.index] = special.off
			}
		}
	}
	return string(text)
}

func (r *runner) timestamp(date time.Time) string {
	// GNU ls uses half the average Gregorian year to select recent timestamps.
	recent := !date.Before(r.now.Add(-15778476*time.Second)) && !date.After(r.now)
	switch r.timeStyle {
	case "full-iso":
		return date.Format("2006-01-02 15:04:05.000000000 -0700")
	case "long-iso":
		return date.Format("2006-01-02 15:04")
	case "iso":
		if recent {
			return date.Format("01-02 15:04")
		}
		return date.Format("2006-01-02 ")
	default:
		if recent {
			return date.Format("Jan _2 15:04")
		}
		return date.Format("Jan _2  2006")
	}
}

func (r *runner) metadata(e entry) (row, error) {
	stat, ok := e.info.Sys().(*syscall.Stat_t)
	if !ok {
		return row{}, fmt.Errorf("Linux file metadata is unavailable")
	}
	result := row{entry: e, mode: permissions(stat.Mode), links: stat.Nlink,
		owner: r.owner(stat.Uid), group: r.group(stat.Gid), blocks: stat.Blocks,
		size: strconv.FormatInt(e.info.Size(), 10), date: r.timestamp(e.info.ModTime())}
	if e.info.Mode()&os.ModeDevice != 0 {
		// Linux's device encoding includes both low and high major/minor bits.
		device := uint64(stat.Rdev)
		major := (device>>8)&0xfff | (device>>32)&0xfffff000
		minor := device&0xff | (device>>12)&0xffffff00
		result.major = strconv.FormatUint(major, 10)
		result.minor = strconv.FormatUint(minor, 10)
	}
	return result, nil
}

// Escape control characters so filenames cannot inject terminal commands or extra lines.
func displayName(name string) string {
	for _, b := range []byte(name) {
		if b < 32 || b == 127 {
			return strconv.Quote(name)
		}
	}
	return name
}

func (r *runner) name(e entry) string {
	name := displayName(e.name)
	if !r.color {
		return name
	}
	color := ""
	switch mode := e.info.Mode(); {
	case mode&os.ModeSymlink != 0:
		color = "36"
	case mode.IsDir():
		color = "34"
	case mode&(os.ModeDevice|os.ModeNamedPipe|os.ModeSocket) != 0:
		color = "33"
	case mode.Perm()&0111 != 0:
		color = "32"
	}
	if color == "" {
		return name
	}
	return "\033[" + color + "m" + name + "\033[0m"
}

func (r *runner) printEntries(entries []entry, total bool, widthExtras []entry) {
	if !r.opts.long {
		for _, e := range entries {
			r.write("%s\n", r.name(e))
		}
		return
	}
	var rows []row
	var blocks int64
	linkWidth, ownerWidth, groupWidth, sizeWidth, majorWidth, minorWidth := 1, 1, 1, 1, 0, 0
	// GNU ls includes directory operands when aligning the initial file listing.
	for i, e := range append(entries, widthExtras...) {
		row, err := r.metadata(e)
		if err != nil {
			r.report(e.path, err, 1)
			continue
		}
		if i < len(entries) {
			rows = append(rows, row)
		}
		blocks += row.blocks
		linkWidth = max(linkWidth, len(strconv.FormatUint(row.links, 10)))
		ownerWidth = max(ownerWidth, len(row.owner))
		groupWidth = max(groupWidth, len(row.group))
		if row.major == "" {
			sizeWidth = max(sizeWidth, len(row.size))
		}
		majorWidth = max(majorWidth, len(row.major))
		minorWidth = max(minorWidth, len(row.minor))
	}
	if majorWidth > 0 {
		sizeWidth = max(sizeWidth, majorWidth+minorWidth+2)
	}
	if total {
		r.write("total %d\n", (blocks+1)/2)
	}
	for _, row := range rows {
		size := row.size
		if row.major != "" {
			size = fmt.Sprintf("%*s, %*s", majorWidth, row.major, minorWidth, row.minor)
		}
		name := r.name(row.entry)
		if row.entry.info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(row.entry.path)
			if err != nil {
				r.report(row.entry.path, err, 1)
			} else {
				name += " -> " + displayName(target)
			}
		}
		r.write("%s %*d %-*s %-*s %*s %s %s\n", row.mode, linkWidth, row.links,
			ownerWidth, row.owner, groupWidth, row.group, sizeWidth, size, row.date, name)
	}
}

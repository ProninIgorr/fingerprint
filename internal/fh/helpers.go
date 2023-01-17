package fh

import (
	"fmt"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ResolvePath {"" -> ".", "~..." -> "user.HomeDir..."} -> Abs
func ResolvePath(path string, usr *user.User) (string, error) {
	var err error
	if path == "" {
		path = "."
	}
	if strings.HasPrefix(path, "~") {
		if usr == nil {
			if userName := os.Getenv("SUDO_USER"); userName != "" { // os.UserHomeDir doesn't work with sudo ... os.Getuid() == 0
				usr, err = user.Lookup(userName)
			} else {
				usr, err = user.Current()
			}
			if err != nil {
				return "", fmt.Errorf("resolving path [%s] failed due to inability to get user info: %w", path, err)
			}
		}
		path = usr.HomeDir + path[1:]
	}
	return filepath.Abs(path)
}

func SafeParentResolvePath(path string, usr *user.User, perm os.FileMode) (string, error) {
	fullPath, err := ResolvePath(path, usr)
	if err != nil {
		return path, err
	}
	dir := path
	if !strings.HasSuffix(path, string(filepath.Separator)) {
		dir = filepath.Dir(fullPath)
	}
	err = os.MkdirAll(dir, perm)
	if err != nil {
		return path, err
	}
	return fullPath, nil
}

// IsDirectory checks whether path is directory and exists
func IsDirectory(path string) (b bool, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, err
		}
	}
	if !fi.IsDir() {
		return false, fmt.Errorf(`not a directory: %v`, path)
	}
	return true, nil
}

// ExpandPatternLists expands {comma separated list} of one level deep
func ExpandPatternLists(pattern string) (result []string, err error) {
	// ...{...}...{...}... -> ... | ...}... | ...}...
	result = strings.Split(pattern, "{")
	if len(result) == 1 { // { - none
		return
	}
	lexs := [][]string{{result[0]}}
	for _, cb := range result[1:] { // { >= 1
		// li1,li2...}s -> li1,li2... | s
		pair := strings.Split(cb, "}") // } - for each chunk splitted by { should exist one } => len(ces) == 2
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid {...} usage in pattern [%s]", pattern)
		}
		var lex []string
		// l1,l2,... -> l1 | l2 | ...
		ls := strings.Split(pair[0], ",")
		for _, s := range ls {
			lex = append(lex, s+pair[1])
		}
		// l1s | l2s | ...
		lexs = append(lexs, lex)
	}
	result = lexs[0]
	for _, lex := range lexs[1:] {
		curAcc := make([]string, 0, len(result)*len(lex))
		for _, s := range lex {
			for _, ac := range result {
				curAcc = append(curAcc, ac+s)
			}
		}
		result = curAcc
	}
	return
}

// BytesToHuman converts 1024 to '1 KiB' etc
func BytesToHuman(src uint64) string {
	if src < 10 {
		return fmt.Sprintf("%d B", src)
	}

	s := float64(src)
	base := float64(1024)
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	e := math.Floor(math.Log(s) / math.Log(base))
	suffix := sizes[int(e)]
	val := math.Floor(s/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// yes, I like Python)
//def sizeof_fmt(num, suffix='B'):
//	for unit in ['','Ki','Mi','Gi','Ti','Pi','Ei','Zi']:
//		if abs(num) < 1024.0:
//			return "%3.1f%s%s" % (num, unit, suffix)
//		num /= 1024.0
//	return "%.1f%s%s" % (num, 'Yi', suffix)

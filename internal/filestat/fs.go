// Package filestat defines and implements basic file statistics interface
package filestat

import (
	"fmt"
	"os"
	"time"
)

type Inode uint64

// FileStat describes a file (use GetFileStat)
type FileStat interface {
	// Path - full path (abs cleaned resolved)
	Path() string
	// BaseName - file name wo dirs
	BaseName() string
	// Inode - Unix inode (or analogue) used here to resolve multiple links to the same file content
	Inode() Inode
	// Size - content size
	Size() int64
	// ModTime - modification time
	ModTime() time.Time
	// String - pretty string for view ...
	String() string
}

// GetFileStat - FileStat builder function (uses os specific func newFileStat)
func GetFileStat(path string) (FileStat, error) {
	if fileInfo, err := os.Lstat(path); err == nil {
		//if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		//	if SymLinkEnabled {
		//		targetPath, err := filepath.EvalSymlinks(path)
		//		if err != nil {
		//			return nil, fmt.Errorf("unresolved symlink [%s]: %w", path, err)
		//		}
		//		targetInfo, err := os.Stat(targetPath)
		//		if err != nil {
		//			return nil, fmt.Errorf("getting target file [%s] Stat for symlink [%s] failed: %w", targetPath, path, err)
		//		}
		//		sfs, err := newFileStat(path, fileInfo, nil, priorFunc, nil)
		//		if err != nil {
		//			return nil, fmt.Errorf("getting FileStat for symlink [%s] failed: %w", path, err)
		//		}
		//		return newFileStat(targetPath, targetInfo, metaKeyFunc, priorFunc, sfs)
		//	} else {
		//		return nil, fmt.Errorf("symlink processing is disabled [%s]", path) // info
		//	}
		//} else {
		return newFileStat(path, fileInfo)
		//}
	} else {
		// errors.Is(err, os.ErrNotExist)
		return nil, fmt.Errorf("getting FileStat for file [%s] failed: %w", path, err)
	}
}

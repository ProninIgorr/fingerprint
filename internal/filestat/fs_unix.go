//go:build aix || darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package filestat

import (
	"io/fs"
	"os"
	"syscall"
	"time"
)

// fileStat implements FileStat interface for *nix os
type fileStat struct {
	path     string
	fileInfo os.FileInfo
	sys      *syscall.Stat_t // *nix specific
}

func (fs *fileStat) String() string {
	return fs.path
}

func (fs *fileStat) Path() string { return fs.path }

func (fs *fileStat) Inode() Inode { return Inode(fs.sys.Ino) }

func (fs *fileStat) Size() int64 { return fs.fileInfo.Size() }

func (fs *fileStat) Perm() fs.FileMode { return fs.fileInfo.Mode().Perm() }

func (fs *fileStat) IsRegular() bool { return fs.fileInfo.Mode().IsRegular() }

func (fs *fileStat) BaseName() string { return fs.fileInfo.Name() }

func (fs *fileStat) ModTime() time.Time { return fs.fileInfo.ModTime() }

func newFileStat(path string, fileInfo os.FileInfo) (FileStat, error) {
	sys := fileInfo.Sys().(*syscall.Stat_t)
	//userOwner, err := user.LookupId(fmt.Sprint(sys.Uid))
	//if err != nil {
	//	return nil, err
	//}
	//groupOwner, err := user.LookupGroupId(fmt.Sprint(sys.Gid))
	//if err != nil {
	//	return nil, err
	//}
	fS := fileStat{
		path:     path,
		fileInfo: fileInfo,
		sys:      sys,
	}
	return &fS, nil
}

// todo: add windows support
// e.g.:
//if runtime.GOOS == "windows" {
//	fileinfo, _ := os.Stat(path)
//	stat := fileinfo.Sys().(*syscall.Win32FileAttributeData)
//	...
//} else {
//	fileinfo, _ := os.Stat(path)
//	stat = fileinfo.Sys().(*syscall.Stat_t)
//	...
//}

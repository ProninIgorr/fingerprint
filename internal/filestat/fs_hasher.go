// // Package fh
// // implements basic hashing functionality
package filestat

//
//import (
//	"crypto/md5"
//	"crypto/sha1"
//	"crypto/sha256"
//	"crypto/sha512"
//	"encoding/hex"
//	"fmt"
//	"hash"
//	"io"
//	"log"
//	"os"
//	"strings"
//)
//
//const (
//	Idle   = ""
//	MD5    = "md5"
//	SHA1   = "sha1"
//	SHA256 = "sha256"
//	SHA512 = "sha512"
//)
//
//const EMPTY_CHECKSUM = "-"
//
//// HashFileFunc - type specifies signature of function that calculates hash / checksum for file [path]
//// initial settings (algo, size, etc.) are taken from closure - see GetHashFileFunc
//type HashFileFunc func(fs FileStat, prefix string) (result string, written int64, err error)
//
//// GetHashFileFunc customizes hasher func
//func GetHashFileFunc(algo string, ndMaxSize int64, inBlocks bool) (HashFileFunc, error) {
//	var fileHasher func() hash.Hash
//	switch strings.ToLower(algo) {
//	case Idle:
//		fileHasher = func() hash.Hash { return &idleHasher{} }
//	case MD5:
//		fileHasher = md5.New
//	case SHA1:
//		fileHasher = sha1.New
//	case SHA256:
//		fileHasher = sha256.New
//	case SHA512:
//		fileHasher = sha512.New
//	default:
//		return nil, fmt.Errorf("invalid value for algo hashing: [%s] - not supported", algo)
//		// fileHasher = func() hash.Hash { return &idleHasher{} }
//		// fileHasher = md5.New
//	}
//	return func(fs FileStat, prefix string) (result string, written int64, err error) {
//		var (
//			h        = fileHasher()
//			size     = fs.Size()
//			dMaxSize = ndMaxSize
//		)
//		if _, ok := h.(*idleHasher); !ok {
//			if inBlocks { // dSize is size in blocks
//				dMaxSize = dMaxSize * fs.Blksize() // files can have different block sizes
//			}
//			file, err := os.Open(fs.Path())
//			if err != nil {
//				return result, written, fmt.Errorf("hasing file [%s] failed: %w", fs.Path(), err)
//			}
//			defer func() {
//				if e := file.Close(); e != nil && err == nil {
//					log.Printf("error while closing file [%s]: %v", fs.Path(), e)
//				}
//			}()
//			switch {
//			case dMaxSize > 0:
//				if size > dMaxSize {
//					size = dMaxSize
//				}
//			case dMaxSize < 0:
//				size = -dMaxSize
//				if size > fs.Size() {
//					size = fs.Size()
//				}
//				if ret, err := file.Seek(-size, io.SeekEnd); err != nil {
//					return result, written, fmt.Errorf("seek file %s (%d) at offset %d is failed with ret = %d: %w", fs.Path(), fs.Size(), -size, ret, err)
//				}
//			case dMaxSize == 0:
//				// size = fs.Size()
//			}
//			if written, err = io.CopyN(h, file, size); err != nil {
//				return result, written, fmt.Errorf("hashing file [%s] is failed - written %d: %w", fs.Path(), written, err)
//			}
//		}
//		checksum := hex.EncodeToString(h.Sum(nil))
//		return fmt.Sprintf("%s:%d:%s:%s", prefix, size, algo, checksum), written, nil
//	}, nil
//}

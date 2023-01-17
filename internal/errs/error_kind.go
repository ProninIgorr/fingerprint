package errs

type Kind uint32

// todo: make kind values as flags so that masks can be used to classify compound errors
// const KindOther Kind = 1 << (32 - 1 - iota)
const (
	KindOther           Kind = iota // Unclassified error. This value is not printed in the error message.
	KindTransient                   // Transient error  todo: use prev Error values
	KindInterrupted                 // Interrupted ( some kind of inconsistency )
	KindInvalidValue                // Invalid value for this type of item.
	KindIO                          // External I/O error such as network failure.
	KindOSOpenFile                  // os.Open errors
	KindOSStat                      // error returned from os.Lstat, os.Stat
	KindFileStat                    // FileStat creation failed
	KindPermission                  // Permission denied.
	KindExist                       // Item already exists.
	KindNotExist                    // Item does not exist.
	KindIsDir                       // Item is a directory.
	KindNotDir                      // Item is not a directory.
	KindFileSystemOther             // Other file system related error.
	KindBrokenLink                  // Link target does not exist.
	KindInternal                    // Internal error (for current errs pipeline impl this kind should be last in this list so that len(Kinds) = int(errs.KindInternal))
)

func (k Kind) String() string {
	switch k {
	case KindOther:
		return "other"
	case KindInvalidValue:
		return "invalid value"
	case KindFileSystemOther:
		return "other file system related"
	case KindInterrupted:
		return "interrupted"
	case KindOSStat:
		return "Stat failed"
	case KindFileStat:
		return "FileStat creation failed"
	case KindPermission:
		return "permission denied"
	case KindIO:
		return "I/O"
	case KindExist:
		return "item already exists"
	case KindNotExist:
		return "item does not exist"
	case KindBrokenLink:
		return "link target does not exist"
	case KindIsDir:
		return "item is a directory"
	case KindNotDir:
		return "item is not a directory"
	case KindInternal:
		return "internal"
	case KindTransient:
		return "transient"
	}
	return "unknown"
}

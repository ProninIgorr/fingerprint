package errs

import (
	"fmt"
	"runtime"
	"strings"
)

var withFrames bool

func WithFrames(b bool) {
	withFrames = b
}

// Frame is a single step in stack trace.
type Frame struct {
	// Func contains a function name.
	Func string
	// Line contains a line number.
	Line int
	// Path contains a file path.
	Path string
}

// String formats Frame to string.
func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s()", f.Path, f.Line, f.Func)
}

type Frames []Frame

func (fs Frames) String() string {
	if len(fs) == 0 {
		return ""
	}
	sb := strings.Builder{}
	for i, fr := range fs {
		sb.WriteString(fmt.Sprintf("-> %d: %s", i, fr))
	}
	return sb.String()
}

// ... or debug.Stack()
func Trace(skip int) (result Frames) {
	if withFrames {
		for {
			pc, path, line, ok := runtime.Caller(skip)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc)
			result = append(result, Frame{
				Func: fn.Name(),
				Line: line,
				Path: path,
			})
			skip++
		}
	}
	return
}

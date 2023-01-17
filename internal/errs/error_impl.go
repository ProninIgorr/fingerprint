package errs

import (
	"time"

	cu "github.com/ProninIgorr/fingerprint/internal/contexts"
)

type errorData struct {
	err      error
	severity Severity
	kind     Kind
	ops      cu.Operations
	frames   []Frame
	ts       int64
}

var _ Error = &errorData{}
var _ Error = (*errorData)(nil)

func newError() Error {
	return errorData{kind: KindOther, severity: SeverityError, ts: time.Now().UnixNano()}
}

func (e errorData) TimeStamp() time.Time {
	return time.Unix(0, e.ts)
}

func (e errorData) Error() string {
	if e.kind == KindTransient {
		if ee, ok := e.err.(Error); ok {
			return ee.Unwrap().Error()
		}
	}
	return e.err.Error()
}

func (e errorData) Severity() Severity {
	return e.severity
}

func (e errorData) Kind() Kind {
	return e.kind
}

func (e errorData) OperationPath() cu.Operations {
	return e.ops
}

func (e errorData) StackTrace() []Frame {
	return e.frames
}

func (e errorData) Unwrap() error {
	return e.err
}

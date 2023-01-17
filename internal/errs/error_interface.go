package errs

import (
	"time"

	cou "github.com/ProninIgorr/fingerprint/internal/contexts"
)

type Error interface {
	error
	Severity() Severity
	TimeStamp() time.Time
	Kind() Kind
	OperationPath() cou.Operations
	StackTrace() []Frame
	Unwrap() error
}

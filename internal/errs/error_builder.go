package errs

import (
	"context"
	"errors"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
)

func E(args ...interface{}) Error {
	switch len(args) {
	case 0:
		panic("call to errors.E with no arguments")
	case 1:
		if e, ok := args[0].(Error); ok {
			return e
		}
	}
	e := newError().(errorData)
	// the last on the list [args] wins
	for _, arg := range args {
		switch a := arg.(type) {
		case Severity:
			e.severity = a
		case Kind:
			e.kind = a
		case contexts.Operation:
			e.ops = contexts.Operations{Stack: []contexts.Operation{a}}
		case contexts.Operations:
			e.ops = a
		case context.Context:
			e.ops = contexts.GetContextOperations(a)
			e.kind = GetDefaultErrsKind(a)
			e.severity = GetDefaultErrsSeverity(a)
		case Error: // todo: impl transient error in this case (need more cases...)
			e.err = a
			if e.kind == KindOther {
				e.kind = KindTransient
			}
		case error:
			e.err = a
		case string:
			e.err = errors.New(a)
		default:
		}
	}
	e.frames = Trace(2)
	return e
}

package errs

import (
	"context"

	cu "github.com/ProninIgorr/fingerprint/internal/contexts"
)

type ctxErrsSeverityKey int

const DefaultErrsSeverityKey ctxErrsSeverityKey = 0

func SetDefaultErrsSeverity(severity Severity) cu.PartialContextFn {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, DefaultErrsSeverityKey, severity)
	}
}

func GetDefaultErrsSeverity(ctx context.Context) Severity {
	if severity, ok := ctx.Value(DefaultErrsSeverityKey).(Severity); ok {
		return severity
	}
	return SeverityError
}

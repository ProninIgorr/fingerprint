package errs

import (
	"context"

	cu "github.com/ProninIgorr/fingerprint/internal/contexts"
)

type ctxErrsKindKey int

const DefaultErrsKindKey = 0

func SetDefaultErrsKind(kind Kind) cu.PartialContextFn {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, DefaultErrsKindKey, kind)
	}
}

func GetDefaultErrsKind(ctx context.Context) Kind {
	if kind, ok := ctx.Value(DefaultErrsKindKey).(Kind); ok {
		return kind
	}
	return KindOther
}

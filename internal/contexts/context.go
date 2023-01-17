package contexts

import (
	"context"
)

func BuildContext(ctx context.Context, ctxFns ...PartialContextFn) context.Context {
	if ctx == nil {
		ctx = context.TODO()
	}
	for _, f := range ctxFns {
		ctx = f(ctx)
	}
	return ctx
}

type PartialContextFn func(context.Context) context.Context

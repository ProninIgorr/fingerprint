package contexts

import (
	"context"
)

type ctxOperationKey int

const OperationKey ctxOperationKey = 0

func AddContextOperation(op Operation) PartialContextFn {
	return func(ctx context.Context) context.Context {
		ops, _ := ctx.Value(OperationKey).(Operations)
		ops.Add(op)
		return context.WithValue(ctx, OperationKey, ops)
	}
}

func SetContextOperation(op Operation) PartialContextFn {
	return func(ctx context.Context) context.Context {
		var ops Operations
		ops.Add(op)
		return context.WithValue(ctx, OperationKey, ops)
	}
}

func GetContextOperations(ctx context.Context) Operations {
	ops, _ := ctx.Value(OperationKey).(Operations)
	return ops
}

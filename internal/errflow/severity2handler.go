package errflow

import (
	"context"
	"sync"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
)

type FuncErrorHandler func(cerr <-chan errs.Error, wg *sync.WaitGroup)

func MapErrorHandlers(
	ctx context.Context,
	scerr map[errs.Severity]chan errs.Error,
	handlers map[errs.Severity]FuncErrorHandler,
	defaultHandler FuncErrorHandler,
) <-chan struct{} {
	ctx = contexts.BuildContext(ctx, contexts.AddContextOperation("handling"))
	done := make(chan struct{})
	var wg sync.WaitGroup
	for severity, cerr := range scerr {
		handler := defaultHandler
		if handlers != nil {
			if _, ok := handlers[severity]; ok {
				handler = handlers[severity]
			}
		}
		wg.Add(1)
		go handler(cerr, &wg)
		logging.LogMsg(ctx).Debug("Errors handlers for [", severity, "] - started")
	}
	go func() {
		wg.Wait()
		close(done)
		logging.LogMsg(ctx).Debug("Errors handlers - stopped")
	}()
	return done
}

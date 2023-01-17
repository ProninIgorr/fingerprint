package workflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/ProninIgorr/fingerprint/internal/errs"
	logging "github.com/ProninIgorr/fingerprint/internal/logging"
)

func OnExit(ctx context.Context, cerr chan<- errs.Error, prefixMsg string, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = errors.New(fmt.Sprintf("%v", r))
			}
			logging.LogError(errs.E(ctx, errs.KindInternal, errs.SeverityCritical, fmt.Errorf("%s - interrupted: %w", prefixMsg, err)))
		}
		fn()
	}()
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			err = errors.New(fmt.Sprintf("%v", r))
		}
		cerr <- errs.E(ctx, errs.KindInternal, errs.SeverityCritical, fmt.Errorf("%s - interrupted: %w", prefixMsg, err))
	} else {
		select {
		case <-ctx.Done():
			cerr <- errs.E(ctx, errs.KindInterrupted, fmt.Errorf("%s - interrupted: %w", prefixMsg, ctx.Err()))
		default:
			logging.LogMsg(ctx).Debug(prefixMsg, " - ok")
		}
	}
}

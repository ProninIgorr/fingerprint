package errflow

import (
	"context"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
	"github.com/ProninIgorr/fingerprint/internal/registrator"
)

type ErrorProducer interface {
	ErrCh() <-chan errs.Error
}

type ErrorStats registrator.Encounter

type ErrorModerator interface {
	Run(ctx context.Context) <-chan struct{}
	Stats() interface{}
}

type errorModerator struct {
	done  <-chan struct{}
	stats ErrorStats
}

func NewErrorModerator(ctx context.Context, cancel context.CancelFunc, errProducers ...ErrorProducer) (ErrorModerator, errs.Error) {
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("_.errs_moderation"))
	errChs := make([]<-chan errs.Error, 0, len(errProducers))
	for _, errProducer := range errProducers {
		errChs = append(errChs, errProducer.ErrCh())
	}
	errsCh := MergeErrors(ctx, errChs...)
	mapSeverity2ErrorCh, totalErrorStats := SortFilteredErrors(ctx, errsCh, logging.GetSeveritiesFilter4CurrentLogLevel())
	errsDone := MapErrorHandlers(
		ctx,
		mapSeverity2ErrorCh,
		nil, // map[errs.Severity]FuncErrorHandler{
		//	errs.SeverityCritical: erf.CriticalErrorHandlerBuilder(cancel, []errs.Kind{...}),
		//}
		LoggingErrorHandler,
	)
	return &errorModerator{
		done:  errsDone,
		stats: totalErrorStats,
	}, nil
}

func (r *errorModerator) Run(context.Context) <-chan struct{} {
	return r.done
}

func (r *errorModerator) Stats() interface{} {
	return r.stats
}

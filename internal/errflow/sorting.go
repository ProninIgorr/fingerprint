package errflow

import (
	"context"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
	"github.com/ProninIgorr/fingerprint/internal/registrator"
)

func SortFilteredErrors(ctx context.Context, inputErrCh <-chan errs.Error, filterSeverities []errs.Severity) (map[errs.Severity]chan errs.Error, ErrorStats) {
	ctx = contexts.BuildContext(ctx, contexts.AddContextOperation("sorting"))
	scerr := make(map[errs.Severity]chan errs.Error)
	stats := ErrorStats(registrator.NewEncounter(len(errs.AllSeverities) * int(errs.KindInternal) * 64))
	for _, severity := range filterSeverities {
		scerr[severity] = make(chan errs.Error, cap(inputErrCh))
	}
	go func() {
		defer func() {
			for severity, cerr := range scerr {
				close(cerr)
				logging.LogMsg(ctx).Debug("Error channel - ", severity.String(), " - closed")
			}
		}()
		for err := range inputErrCh {
			if err != nil {
				if cerr, ok := scerr[err.Severity()]; ok {
					cerr <- err
				}
				stats.CheckIn(ErrStatKey{err.Severity(), err.Kind(), err.OperationPath().String()})
			}
		}
	}()
	return scerr, stats
}

package errflow

import (
	"context"
	"sync"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
)

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func MergeErrors(ctx context.Context, errChs ...<-chan errs.Error) <-chan errs.Error {
	ctx = contexts.BuildContext(ctx, contexts.AddContextOperation("merging"))
	var wg sync.WaitGroup
	// We must ensure that the output channel has the reading capacity to hold as many errors
	// as there could be written to all error channels at once.
	// This will ensure that it never blocks, even
	// if further processing ended before closing the channel.
	var capOut int
	for _, errCh := range errChs {
		if errCh != nil {
			capOut += cap(errCh)
		}
	}
	outputErrCh := make(chan errs.Error, capOut)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(errCh <-chan errs.Error) {
		for err := range errCh {
			outputErrCh <- err
		}
		wg.Done()
	}
	wg.Add(len(errChs))
	for _, errCh := range errChs {
		if errCh != nil {
			go output(errCh)
		}
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(outputErrCh)
		logging.LogMsg(ctx).Debug("Merged errors channel - closed")
	}()
	return outputErrCh
}

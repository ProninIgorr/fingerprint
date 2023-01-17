package inputs

import (
	"context"
	"fmt"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/logging"
	"github.com/ProninIgorr/fingerprint/internal/registrator"
	"github.com/ProninIgorr/fingerprint/internal/workflow"
	"github.com/fsnotify/fsnotify"
)

type InputsStats registrator.Encounter

type Inputs interface {
	FoundFilePathsCh() <-chan string
	ErrCh() <-chan errs.Error
	Run(ctx context.Context) <-chan struct{}
	Stats() interface{}
}

type inputs struct {
	dir     string
	pattern string
	resCh   chan string
	errCh   chan errs.Error
	stats   InputsStats
}

func NewInputs(ctx context.Context, dir string, filePattern string, initCap int) Inputs {
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("1.0.inputs_init"))
	sr := inputs{
		dir:     dir,
		pattern: filePattern,
		resCh:   make(chan string, initCap),
		errCh:   make(chan errs.Error, 2),
		stats:   InputsStats(registrator.NewEncounter(initCap * 8)),
	}
	return &sr
}

func (r *inputs) Run(ctx context.Context) <-chan struct{} {
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("1.inputs_run"))
	done := make(chan struct{})

	go func() {
		ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("inputs_worker"))
		logging.LogMsg(ctx).Debugf(fmt.Sprintf("inputs [%s]:[%s] - started", r.dir, r.pattern))
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			r.errCh <- errs.E(ctx, errs.KindIO, errs.SeverityCritical, fmt.Errorf("fsnotify.NewWatcher() failed: %w", err))
			return
		}
		defer workflow.OnExit(ctx, r.errCh, "workers", func() {
			close(r.resCh)
			close(r.errCh)
			close(done)
			watcher.Close()
		})
		if err = watcher.Add(r.dir); err != nil {
			r.errCh <- errs.E(ctx, errs.KindIO, errs.SeverityCritical, fmt.Errorf("adding directory [%s] for watching failed: %w", r.dir, err))
			return
		}
		logging.LogMsg(ctx).Debugf(fmt.Sprintf("inputs [%s]:[%s] - started", r.dir, r.pattern))

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				logging.LogMsg(ctx).Debugf(fmt.Sprintf("event [%s] - fired with ok = %v", event, ok))
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					path := event.Name
					logging.LogMsg(ctx).Debugf(fmt.Sprintf("file [%s] - added", path))
					r.stats.CheckIn(path)
					r.resCh <- path
				}
			case err, ok := <-watcher.Errors:
				r.errCh <- errs.E(ctx, errs.KindIO, fmt.Errorf("watcher error with ok = [%v]: %w", ok, err))
				if !ok {
					return
				}
			}
		}
	}()
	return done
}

func (r *inputs) FoundFilePathsCh() <-chan string {
	return r.resCh
}

func (r *inputs) ErrCh() <-chan errs.Error {
	return r.errCh
}

func (r *inputs) Stats() interface{} {
	return r.stats
}

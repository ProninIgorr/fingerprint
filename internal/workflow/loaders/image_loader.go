package loaders

import (
	"context"
	"fmt"
	"sync"

	"github.com/ProninIgorr/fingerprint/internal/contexts"
	"github.com/ProninIgorr/fingerprint/internal/errs"
	"github.com/ProninIgorr/fingerprint/internal/filestat"
	"github.com/ProninIgorr/fingerprint/internal/imgs/matrix"
	"github.com/ProninIgorr/fingerprint/internal/registrator"
	"github.com/ProninIgorr/fingerprint/internal/workflow"
)

type LoaderStats struct {
	FileStats registrator.Encounter
	//InodeStats registrator.Encounter
}

type ImageLoader interface {
	ImageCh() <-chan workflow.ImageStat
	ErrCh() <-chan errs.Error
	Run(ctx context.Context) <-chan struct{}
	Stats() interface{}
}

type LoadImageFunc func(fname string, maxDimension int) (*matrix.M, error)

type imageLoader struct {
	inputCh <-chan filestat.FileStat
	resCh   chan workflow.ImageStat
	errCh   chan errs.Error
	stats   LoaderStats

	fn         LoadImageFunc
	maxWorkers int
}

func NewImageLoader(ctx context.Context,
	inputCh <-chan filestat.FileStat,
	fn LoadImageFunc,
	initCap int) ImageLoader {
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("2.0.image_loader_init"))
	maxWorkers := cap(inputCh) // * runtime.NumCPU()
	loader := imageLoader{
		inputCh: inputCh,
		resCh:   make(chan workflow.ImageStat, maxWorkers),
		errCh:   make(chan errs.Error, maxWorkers*2),
		stats: LoaderStats{
			FileStats: registrator.NewEncounter(initCap),
			//InodeStats: registrator.NewEncounter(initCap),
		},
		fn:         fn,
		maxWorkers: maxWorkers,
	}
	return &loader
}

func (r *imageLoader) Run(ctx context.Context) <-chan struct{} {
	ctx = contexts.BuildContext(ctx, contexts.SetContextOperation("image_loader_run"))
	done := make(chan struct{})

	go func() {
		ctx = contexts.BuildContext(ctx, contexts.AddContextOperation("image_loader_workers"))
		var (
			wg    sync.WaitGroup
			wPool = make(chan struct{}, r.maxWorkers)
		)
		defer workflow.OnExit(ctx, r.errCh, "workers", func() {
			wg.Wait()
			close(wPool)
			close(r.resCh)
			close(r.errCh)
			close(done)
		})
		for {
			select {
			case <-ctx.Done():
				return
			case fs, more := <-r.inputCh:
				if !more {
					return
				}
				select {
				case <-ctx.Done():
					return
				case wPool <- struct{}{}:
					wg.Add(1)
					go func(filestat.FileStat) {
						defer wg.Done()
						defer func() { <-wPool }()
						if m, err := r.fn(fs.Path(), 300); err == nil {
							imageStat := workflow.ImageStat{
								Fs: fs,
								IM: m,
							}
							select {
							case <-ctx.Done():
								return
							case r.resCh <- imageStat:
								r.stats.FileStats.CheckIn(registrator.KeySize{Key: fs.String(), Size: fs.Size()})
							}
						} else {
							r.errCh <- errs.E(ctx, errs.KindFileStat, fmt.Errorf("loading [%s] failed: %w", fs.Path(), err))
						}
					}(fs)
				}
			}
		}
	}()
	return done
}

func (r *imageLoader) ImageCh() <-chan workflow.ImageStat {
	return r.resCh
}

func (r *imageLoader) ErrCh() <-chan errs.Error {
	return r.errCh
}

func (r *imageLoader) Stats() interface{} {
	return &r.stats
}

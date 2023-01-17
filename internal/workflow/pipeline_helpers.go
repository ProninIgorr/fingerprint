package workflow

import (
	"context"

	"github.com/ProninIgorr/fingerprint/internal/logging"
)

type Runner interface {
	Run(context.Context) <-chan struct{}
}

func Run(ctx context.Context, runners ...Runner) <-chan struct{} {
	doneChs := make([]<-chan struct{}, 0, len(runners))
	for _, runner := range runners {
		doneChs = append(doneChs, runner.Run(ctx))
	}
	done := make(chan struct{})
	go func() {
		for i, dc := range doneChs {
			if dc != nil {
				logging.LogMsg(ctx).Debugf("closing %T", runners[i])
				<-dc
			}
		}
		close(done)
	}()
	return done
}

type StatProducer interface {
	Stats() interface{}
}

type Pipeliner interface {
	Runner
	StatProducer
}

type Pipelines []Pipeliner

func (pl Pipelines) Runners() []Runner {
	result := make([]Runner, 0, len(pl))
	for _, p := range pl {
		result = append(result, p)
	}
	return result
}

func (pl Pipelines) StatProducers() []StatProducer {
	result := make([]StatProducer, 0, len(pl))
	for _, p := range pl {
		result = append(result, p)
	}
	return result
}

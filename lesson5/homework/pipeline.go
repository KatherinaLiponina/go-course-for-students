package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	for i := 0; i < len(stages); i++ {
		in = stages[i](in);
	}
	out := make(chan any)
	go func() {
		for {
			select {
			case r, ok := <-in:
				if !ok {
					close(out)
					return
				}
				out <- r
			case <-ctx.Done():
				close(out)
				return
			}
		}
	}()
	return out
}

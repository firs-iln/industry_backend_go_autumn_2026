package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	results := make([]Result[R], len(in))
	if len(in) == 0 {
		return results, nil
	}
	// One goroutine per item is not allowed, and idle workers are pointless.
	if workers > len(in) {
		workers = len(in)
	}
	// Every index is queued up front, so a worker that finishes early picks up
	// the next item itself and the pool stays busy while work is left.
	jobs := make(chan int, len(in))
	for i := range in {
		jobs <- i
	}
	close(jobs)

	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			for i := range jobs {
				// Cancellation stops handing out items; a call already in
				// flight keeps running and is waited for below.
				if ctx.Err() != nil {
					return
				}
				value, err := fn(ctx, in[i])
				if err != nil {
					// An item's error is its own: it leaves Value zero and
					// does not stop the other items.
					results[i] = Result[R]{Err: err}
					continue
				}
				results[i] = Result[R]{Value: value}
			}
		})
	}
	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"sync"
	"time"
)

const (
	// DefaultWorkers is the default initial worker count for a Runner.
	DefaultWorkers = 1
	// DefaultMaxWorkers is the default maximum worker count for a Runner.
	DefaultMaxWorkers = 10
)

// Runner is a benchmark runner for a Resolver.
type Runner struct {
	resolver   Resolver
	stop       chan struct{}
	workers    int
	maxWorkers int
	qps        int
	begin      time.Time
	end        time.Time
}

// NewRunner returns a new Runner
func NewRunner(resolver Resolver, opts ...func(*Runner)) *Runner {
	r := &Runner{
		resolver:   resolver,
		stop:       make(chan struct{}),
		workers:    DefaultWorkers,
		maxWorkers: DefaultMaxWorkers,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Workers returns an options function which sets the initial worker count for a Runner.
func Workers(workers int) func(*Runner) {
	return func(r *Runner) { r.workers = workers }
}

// MaxWorkers returns an options function which sets the maximum worker count for a Runner.
func MaxWorkers(workers int) func(*Runner) {
	return func(r *Runner) { r.maxWorkers = workers }
}

// Start initiates a benchmark run using the provided Namer, QPS target, and total request count.
func (r *Runner) Start(names Namer, qps int, count int) <-chan *Result {
	var workers sync.WaitGroup
	res := make(chan *Result)
	run := make(chan int)
	var interval int64
	for i := int(0); i < r.workers; i++ {
		workers.Add(1)
		go r.run(names, &workers, run, res)
	}

	go func() {
		defer close(res)
		defer workers.Wait()
		defer close(run)
		r.begin = time.Now()

		if qps > 0 {
			interval = time.Second.Nanoseconds() / int64(qps)
		}
		requests := int(1)
		for {
			if qps > 0 {
				now := time.Now()
				next := r.begin.Add(time.Duration(interval * int64(requests)))
				time.Sleep(next.Sub(now))
			}
			select {
			case run <- requests:
				if requests++; count > 0 && requests > count {
					return
				}
			case <-r.stop:
				return
			default:
				if r.maxWorkers > 0 && r.workers < r.maxWorkers {
					workers.Add(1)
					go r.run(names, &workers, run, res)
					r.workers++
				}
			}
		}
	}()

	return res
}

// Stop shuts down a running test.
func (r *Runner) Stop() {
	select {
	case <-r.stop:
		break
	default:
		close(r.stop)
	}
	r.end = time.Now()
}

// Begin returns the time the test was started.
func (r *Runner) Begin() time.Time {
	return r.begin
}

// End returns the time the test wast stopped.
func (r *Runner) End() time.Time {
	return r.end
}

// NumWorkers returns the number of active workers.
func (r *Runner) NumWorkers() int {
	return r.workers
}

func (r *Runner) run(name Namer, workers *sync.WaitGroup, run <-chan int, res chan<- *Result) {
	defer workers.Done()
	for range run {
		res <- r.resolver.Resolve(name)
	}
}

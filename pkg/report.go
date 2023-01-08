// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"fmt"
	"io"
	"time"

	"github.com/codahale/hdrhistogram"
)

const ms float64 = 1000000

// TextReporter is a reporter that writes simple text formatted information intented for human reading.
type TextReporter struct {
	interval  time.Duration
	count     int64
	errors    int64
	hist      *hdrhistogram.Histogram
	intErrors int64
	intHist   *hdrhistogram.Histogram
	writer    io.Writer
}

// NewTextReporter returns a new TextReporter
func NewTextReporter(w io.Writer, start time.Time, interval time.Duration) *TextReporter {
	return &TextReporter{
		hist:     hdrhistogram.New(0, 5e+9, 3),
		intHist:  hdrhistogram.New(0, 5e+9, 3),
		interval: interval,
		writer:   w,
	}
}

// Display writes interval statistics to the configured output destination.
func (r *TextReporter) Display() {
	r.display(r.intHist, r.intErrors, float64(r.interval)/float64(time.Second))
}

// Summary writes end of run details.
func (r *TextReporter) Summary(run *Runner) {
	durSeconds := float64(run.End().Sub(run.Begin())) / float64(time.Second)
	fmt.Fprintf(r.writer, "\nFinished %d requests\n\n", r.count)
	fmt.Fprintf(r.writer, "# latency summary\n")
	r.display(r.hist, r.errors, float64(time.Since(run.Begin()))/float64(time.Second))
	fmt.Fprintf(r.writer, "\nConcurrency level: %d\n", run.NumWorkers())
	fmt.Fprintf(r.writer, "Time taken for tests: %s\n", time.Since(run.Begin()))
	fmt.Fprintf(r.writer, "Completed Requests: %d\n", r.count-r.errors)
	fmt.Fprintf(r.writer, "Failed Requests: %d\n", r.errors)
	fmt.Fprintf(r.writer, "Requests per second: %.4f [#/sec] (mean)\n", float64(r.count)/durSeconds)
	fmt.Fprintf(r.writer, "Time per request: %.2f [ms] (mean)\n", float64(r.hist.Mean())/ms)
	fmt.Fprintf(r.writer, "Fastest request: %.2f [ms]\n", float64(r.hist.Min())/ms)
	fmt.Fprintf(r.writer, "Slowest request: %.2f [ms]\n", float64(r.hist.Max())/ms)
}

// Update records a new result.
func (r *TextReporter) Update(res *Result) {
	r.count++
	if res.Error != nil {
		r.errors++
		r.intErrors++
	}
	r.hist.RecordValue(int64(res.Duration))
	r.intHist.RecordValue(int64(res.Duration))
}

// NewInterval resets the current interval statistics.
func (r *TextReporter) NewInterval() {
	r.intErrors = 0
	r.intHist.Reset()
}

// Heading writes the report heading.
func (r *TextReporter) Heading() {
	fmt.Fprintf(r.writer, "# requests errors min  [ p50  p95  p99  p999] max  qps\n")
}

func (r *TextReporter) display(hist *hdrhistogram.Histogram, errors int64, dur float64) {
	fmt.Fprintf(r.writer, "%7d %7d %6.2f    [%4.2f %4.2f %4.2f %5.2f] %3.2f  %3.2f\n",
		hist.TotalCount(),
		errors,
		float64(hist.Min())/ms,
		float64(hist.ValueAtQuantile(50))/ms,
		float64(hist.ValueAtQuantile(95))/ms,
		float64(hist.ValueAtQuantile(99))/ms,
		float64(hist.ValueAtQuantile(999))/ms,
		float64(hist.Max())/ms,
		float64(hist.TotalCount())/dur,
	)
}

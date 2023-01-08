// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestWorkersOpt(t *testing.T) {
	r := NewRunner(nil, Workers(10))
	if r.workers != 10 {
		t.Errorf("got %d, want: 10", r.workers)
	}
}

func TestMaxWorkersOpt(t *testing.T) {
	r := NewRunner(nil, MaxWorkers(100))
	if r.maxWorkers != 100 {
		t.Errorf("got %d, want: 100", r.maxWorkers)
	}
}

func TestRequestCount(t *testing.T) {
	dns.HandleFunc("good.invalid.", handler)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()
	n, err := FileNamer(strings.NewReader("good.invalid"))
	r := NewRunner(NewRemoteResolver(Nameserver(l)))
	requests := 0
	for range r.Start(n, 100, 100) {
		requests++
	}
	r.Stop()
	if requests != 100 {
		t.Errorf("got %v results, want: 100", requests)
	}
}

func TestMaxWorkers(t *testing.T) {
	dns.HandleFunc("good.invalid.", handlerSlow)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()
	n, err := FileNamer(strings.NewReader("good.invalid"))
	r := NewRunner(NewRemoteResolver(Nameserver(l)), MaxWorkers(1))
	requests := 0
	for range r.Start(n, 100, 100) {
		requests++
	}
	r.Stop()
	if requests != 100 {
		t.Errorf("got %v results, want: 100", requests)
	}
	if r.NumWorkers() != 1 {
		t.Errorf("got %v workers, want: 1", r.NumWorkers())
	}

}

func TestTimeout(t *testing.T) {
	dns.HandleFunc("good.invalid.", handlerSlow)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()
	n, err := FileNamer(strings.NewReader("good.invalid"))
	r := NewRunner(NewRemoteResolver(Nameserver(l), Timeout(time.Millisecond)), MaxWorkers(1))
	requests := 0
	var res *Result
	for res = range r.Start(n, 1, 1) {
		requests++
	}
	r.Stop()
	if requests != 1 {
		t.Errorf("got %v, want: 1", requests)
	}
	if res.Error == nil {
		t.Error("got nil, want: error")
	}
	want := "i/o timeout"
	if got := res.Error.Error(); !strings.Contains(got, want) {
		t.Errorf("want: '%v' in '%v'", want, got)
	}
	if res.Error == nil {
		t.Error("got nil, want: error")
	}
}

func TestQPS(t *testing.T) {
	target := 100
	dns.HandleFunc("good.invalid.", handler)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()
	n, err := FileNamer(strings.NewReader("good.invalid"))
	r := NewRunner(NewRemoteResolver(Nameserver(l)))
	requests := 0
	for range r.Start(n, target, target) {
		requests++
	}
	r.Stop()
	if requests != target {
		t.Errorf("got %v results, want: %v", requests, target)
	}
	duration := r.End().Sub(r.Begin())
	qps := int(math.Round(float64(requests) / (float64(duration) / float64(time.Second))))
	if qps != target && qps != target-1 && qps != target+1 {
		t.Errorf("got %v qps, want: %v+/-1", qps, target)
	}
}

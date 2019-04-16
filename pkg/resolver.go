// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

const (
	DefaultNameserver = "8.8.8.8:53"
	DefaultTimeout    = time.Second * 5
)

type Resolver interface {
	Resolve(name Namer) *Result
}

type RemoteResolver struct {
	nameserver string
	timeout    time.Duration
}

func NewRemoteResolver(opts ...func(*RemoteResolver)) *RemoteResolver {
	r := &RemoteResolver{
		nameserver: DefaultNameserver,
		timeout:    DefaultTimeout,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func Timeout(timeout time.Duration) func(*RemoteResolver) {
	return func(r *RemoteResolver) { r.timeout = timeout }
}

func Nameserver(ns string) func(*RemoteResolver) {
	return func(r *RemoteResolver) { r.nameserver = ns }
}

func (r *RemoteResolver) Resolve(name Namer) *Result {
	c := &dns.Client{Timeout: r.timeout}
	m := &dns.Msg{}
	n := name()
	m.SetQuestion(n.Name, n.Type)
	res, dur, err := c.Exchange(m, r.nameserver)
	if err == nil && res.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("lookup failed: %s", dns.RcodeToString[res.Rcode])
	}
	return &Result{
		Duration: dur,
		Error:    err,
	}
}

type LocalResolver struct{}

func NewLocalResolver(opts ...func(*LocalResolver)) *LocalResolver {
	return &LocalResolver{}
}

func (r *LocalResolver) Resolve(name Namer) *Result {
	start := time.Now()
	_, err := net.LookupHost(name().Name)
	return &Result{
		Duration: time.Since(start),
		Error:    err,
	}
}

// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"net"
	"strings"
	"testing"
	"time"

	"github.com/miekg/dns"
)

func TestLocalResolver(t *testing.T) {
	r := NewLocalResolver()
	n, err := FileNamer(strings.NewReader("example.com\ntest.invalid"))
	if err != nil {
		t.Error(err)
	}
	if r.Resolve(n).Error != nil {
		t.Errorf("unexpected error `%s` with valid hostname `example.com`.", err)
	}
	if r.Resolve(n).Error == nil {
		t.Error("expected error with invalid hostname `test.invalid`.")
	}
}

func TestTimeoutOpt(t *testing.T) {
	timeout := time.Second * 100
	r := NewRemoteResolver(Timeout(timeout))
	if r.timeout != timeout {
		t.Errorf("expected timeout to be 100 seconds got %s", r.timeout)
	}
}

func TestNameserverOpt(t *testing.T) {
	ns := "127.0.0.1:53"
	r := NewRemoteResolver(Nameserver(ns))
	if r.nameserver != ns {
		t.Errorf("expected nameserver to be 127.0.0.1:53 got %s", r.nameserver)
	}
}

func TestRemoteResolver(t *testing.T) {
	dns.HandleFunc("good.invalid.", handler)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()

	r := NewRemoteResolver(Nameserver(l))
	n, err := FileNamer(strings.NewReader("good.invalid\nbad.invalid"))
	if err != nil {
		t.Error(err)
	}
	if r.Resolve(n).Error != nil {
		t.Errorf("unexpected error `%s` with valid hostname `good.invalid`.", err)
	}
	if r.Resolve(n).Error == nil {
		t.Error("expected error with invalid hostname `bad.invalid`.")
	}
}

func TestRemoteResolverTimeout(t *testing.T) {
	dns.HandleFunc("good.invalid.", handlerSlow)
	defer dns.HandleRemove("good.invalid.")
	s, l, err := runServer(":0")
	if err != nil {
		t.Error(err)
	}
	defer s.Shutdown()

	r := NewRemoteResolver(Nameserver(l), Timeout(time.Millisecond))
	n, err := FileNamer(strings.NewReader("good.invalid"))
	if err != nil {
		t.Error(err)
	}
	err = r.Resolve(n).Error
	if err == nil {
		t.Error("Expected error")
	}
	ne, ok := err.(*net.OpError)
	if !ok {
		t.Errorf("Expected net.OpError got %T", err)
	}
	if !ne.Timeout() {
		t.Errorf("Expected a timeout error got: %s", ne)
	}
}

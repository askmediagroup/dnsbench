// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
)

func handler(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Answer = make([]dns.RR, 1)
	m.Answer[0] = &dns.A{Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}, A: net.ParseIP("127.0.0.1")}
	w.WriteMsg(m)
}

func handlerSlow(w dns.ResponseWriter, req *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Answer = make([]dns.RR, 1)
	m.Answer[0] = &dns.A{Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0}, A: net.ParseIP("127.0.0.1")}
	time.Sleep(time.Millisecond * 10)
	w.WriteMsg(m)
}

func runServer(laddr string, opts ...func(*dns.Server)) (*dns.Server, string, error) {
	pc, err := net.ListenPacket("udp", laddr)
	if err != nil {
		return nil, "", err
	}
	server := &dns.Server{PacketConn: pc, ReadTimeout: time.Hour, WriteTimeout: time.Hour}

	waitLock := sync.Mutex{}
	waitLock.Lock()
	server.NotifyStartedFunc = waitLock.Unlock

	for _, opt := range opts {
		opt(server)
	}

	go func() {
		server.ActivateAndServe()
		pc.Close()
	}()

	waitLock.Lock()
	return server, pc.LocalAddr().String(), nil
}

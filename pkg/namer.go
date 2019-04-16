// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"bufio"
	"io"
	"sync/atomic"

	"github.com/miekg/dns"
)

type Name struct {
	Name string
	Type uint16
}

type Namer func() Name

func FileNamer(file io.Reader) (Namer, error) {
	var names []Name
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, Name{Name: dns.Fqdn(scanner.Text()), Type: dns.TypeA})
	}
	i := int64(-1)
	n := func() Name {
		return names[atomic.AddInt64(&i, 1)%int64(len(names))]
	}
	return n, scanner.Err()
}

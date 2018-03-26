// Package dnsbench exports types and functions to run dns benchmark tests
package dnsbench

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/codahale/hdrhistogram"
	"github.com/miekg/dns"
)

const ms float64 = 1000000

// Result holds the result of an individual dns lookup
type Result struct {
	Duration time.Duration
	Error    bool
}

// Stats holds a summary of DNS benchmark results
type Stats struct {
	Errors    int64
	Hist      *hdrhistogram.Histogram
	IntErrors int64
	IntHist   *hdrhistogram.Histogram
}

// Resolve does DNS lookups using specified nameserver
func Resolve(nameserver string, names []string, count int, qps int, results chan<- *Result) {
	var result *Result
	var sleep time.Duration
	indexes := make([]int32, count)
	c := new(dns.Client)
	c.Timeout = time.Second * 5
	m := new(dns.Msg)
	for i := range indexes {
		indexes[i] = rand.Int31n(int32(len(names)))
	}
	if qps > 0 {
		sleep = time.Duration(int(time.Second) / qps)
	}

	for i := 0; i < count; i++ {
		m.SetQuestion(dns.Fqdn(names[indexes[i]]), dns.TypeA)
		_, dur, err := c.Exchange(m, net.JoinHostPort(nameserver, "53"))
		result = new(Result)
		result.Duration = dur
		if err != nil {
			result.Error = true
		}
		results <- result
		time.Sleep(sleep)
	}
}

// LocalResolve does DNS lookups using the local resolver
func LocalResolve(names []string, count int, qps int, results chan<- *Result) {
	var result *Result
	var sleep time.Duration
	indexes := make([]int32, count)
	for i := range indexes {
		indexes[i] = rand.Int31n(int32(len(names)))
	}
	if qps > 0 {
		sleep = time.Duration(int(time.Second) / qps)
	}
	for i := 0; i < count; i++ {
		start := time.Now()
		_, err := net.LookupHost(names[indexes[i]])
		result = new(Result)
		result.Duration = time.Since(start)
		if err != nil {
			result.Error = true
		}
		results <- result
		time.Sleep(sleep)
	}
}

// ReadNames reads domain names from a given file and returns them in a slice
func ReadNames(path string) ([]string, error) {
	var file *os.File
	var err error
	if path == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	fmt.Printf("Reading names from %s\n", file.Name())
	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names, scanner.Err()
}

// DisplayReport prints a HDR histogram summary of a DNS benchmark run
func DisplayReport(hist *hdrhistogram.Histogram, errors int64, dur float64) {
	fmt.Printf("%7d %7d %6.2f    [%4.2f %4.2f %4.2f %5.2f] %3.2f  %3.2f\n",
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

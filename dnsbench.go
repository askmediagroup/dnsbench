package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

type Result struct {
	duration time.Duration
	error    bool
}

type Stats struct {
	errors  int
	average int64
	dur     int64
	max     int64
	min     int64
}

func (s *Stats) Update(newRTT time.Duration, n int64) {
	s.average = (int64(newRTT) + (n * int64(s.average))) / (n + 1)
	s.dur += int64(newRTT)
	if int64(newRTT) > s.max {
		s.max = int64(newRTT)
	}
	if s.min == 0 || int64(newRTT) < s.min {
		s.min = int64(newRTT)
	}
}

func resolve(nameserver string, names []string, count int, results chan<- *Result) {
	var result *Result
	indexes := make([]int32, count)
	c := new(dns.Client)
	m := new(dns.Msg)
	for i := range indexes {
		indexes[i] = rand.Int31n(int32(len(names)))
	}

	for i := 0; i < count; i++ {
		m.SetQuestion(dns.Fqdn(names[indexes[i]]), dns.TypeA)
		_, dur, err := c.Exchange(m, net.JoinHostPort(nameserver, "53"))
		result = new(Result)
		result.duration = dur
		if err != nil {
			result.error = true
		}
		results <- result
	}
}

func readNames(path string) ([]string, error) {
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

func main() {
	var count int
	var concur int
	var namesfile string
	var nameserver string
	var result *Result
	stats := Stats{}

	flag.IntVar(&count, "count", 1000, "Number of requests to make")
	flag.IntVar(&concur, "concurency", 10, "Numer of concurrent requests")
	flag.StringVar(&namesfile, "names", "-", "File containing names to request")
	flag.StringVar(&nameserver, "nameserver", "", "Nameserver to query")
	flag.Parse()
	results := make(chan *Result, count)
	names, err := readNames(namesfile)
	if err != nil {
		log.Fatalf("Unable to parse names list: %s", err)
	}

	fmt.Printf("Benchmarking %s...\n\n", nameserver)
	perworker := count / concur
	perworkerrem := count % concur
	start := time.Now().UnixNano()
	for i := 0; i < concur; i++ {
		t := perworker
		if i < perworkerrem {
			t += 1
		}
		go resolve(nameserver, names, t, results)
	}

	mod := 100
	if count >= 1000 {
		mod = count / 10
	}

	for j := 0; j < count; j++ {
		if j > 0 && j%mod == 0 {
			fmt.Printf("Completed %d requests\n", j)
		}
		result = <-results
		stats.Update(result.duration, int64(j))
		if result.error {
			stats.errors++
		}
	}
	end := time.Now().UnixNano()
	duration := float64(end-start) / float64(time.Second)

	fmt.Printf("Finished %d requests\n\n", count)
	fmt.Printf("Concurrency level: %d\n", concur)
	fmt.Printf("Time taken for tests: %.2f seconds\n", duration)
	fmt.Printf("Completed Requests: %d\n", count-stats.errors)
	fmt.Printf("Failed Requests: %d\n", stats.errors)
	fmt.Printf("Requests per second: %.2f [#/sec] (mean)\n", float64(count)/duration)
	fmt.Printf("Time per request: %.2f [ms] (mean)\n", float64(stats.average)/float64(time.Millisecond))
	fmt.Printf("Fastest request: %.2f [ms]\n", float64(stats.min)/float64(time.Millisecond))
	fmt.Printf("Slowest request: %.2f [ms]\n", float64(stats.max)/float64(time.Millisecond))
}

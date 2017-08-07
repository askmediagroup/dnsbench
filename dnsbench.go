package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"time"

	"github.com/codahale/hdrhistogram"
	"github.com/miekg/dns"
)

const ms float64 = 1000000

type Result struct {
	duration time.Duration
	error    bool
}

type Stats struct {
	errors    int64
	hist      *hdrhistogram.Histogram
	intErrors int64
	intHist   *hdrhistogram.Histogram
}

func resolve(nameserver string, names []string, count int, qps int, results chan<- *Result) {
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
		result.duration = dur
		if err != nil {
			result.error = true
		}
		results <- result
		time.Sleep(sleep)
	}
}

func localResolve(names []string, count int, qps int, results chan<- *Result) {
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
		result.duration = time.Since(start)
		if err != nil {
			result.error = true
		}
		results <- result
		time.Sleep(sleep)
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

func displayReport(hist *hdrhistogram.Histogram, errors int64, dur float64) {
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

func main() {
	var nameserver string
	stats := Stats{}
	stats.hist = hdrhistogram.New(0, 5e+9, 3)
	stats.intHist = hdrhistogram.New(0, 5e+9, 3)

	count := flag.Int("count", 1000, "Number of requests to make")
	interval := flag.Duration("interval", 5*time.Second, "Reporting interval")
	concur := flag.Int("concurrency", 10, "Numer of concurrent requests")
	namesfile := flag.String("names", "-", "File containing names to request")
	localResolver := flag.Bool("local", false, "Use local resolver")
	qps := flag.Int("qps", 0, "QPS to attempt for each worker (0 for no limit)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <nameserver> [flags]\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if !*localResolver {
		if flag.NArg() != 1 {
			flag.Usage()
			os.Exit(2)
		}
		nameserver = flag.Arg(0)
	}
	results := make(chan *Result, *concur)
	cleanup := make(chan bool, 2)
	names, err := readNames(*namesfile)
	if err != nil {
		log.Fatalf("Unable to parse names list: %s", err)
	}
	timeout := time.After(*interval)
	if *localResolver {
		fmt.Printf("Benchmarking...\n\n")
	} else {
		fmt.Printf("Benchmarking %s...\n\n", nameserver)
	}
	perworker := *count / *concur
	perworkerrem := *count % *concur
	start := time.Now().UnixNano()
	for i := 0; i < *concur; i++ {
		t := perworker
		if i < perworkerrem {
			t += 1
		}
		if *localResolver {
			go localResolve(names, t, *qps, results)
		} else {
			go resolve(nameserver, names, t, *qps, results)
		}
	}

	fmt.Printf("# requests errors min  [ p50  p95  p99  p999] max  qps\n")
	for {
		select {
		case result := <-results:
			stats.hist.RecordValue(int64(result.duration))
			stats.intHist.RecordValue(int64(result.duration))
			if result.error {
				stats.errors++
				stats.intErrors++
			}
			if stats.hist.TotalCount() >= int64(*count) {
				cleanup <- true
			}
		case <-timeout:
			displayReport(stats.intHist, stats.intErrors, float64(*interval)/float64(time.Second))
			stats.intHist.Reset()
			stats.intErrors = 0
			timeout = time.After(*interval)
		case <-cleanup:
			end := time.Now().UnixNano()
			duration := float64(end-start) / float64(time.Second)

			fmt.Printf("\nFinished %d requests\n\n", *count)
			fmt.Printf("# latency summary\n")
			displayReport(stats.hist, stats.errors, duration)
			fmt.Printf("\nConcurrency level: %d\n", *concur)
			fmt.Printf("Time taken for tests: %.2f seconds\n", duration)
			fmt.Printf("Completed Requests: %d\n", int64(*count)-stats.errors)
			fmt.Printf("Failed Requests: %d\n", stats.errors)
			fmt.Printf("Requests per second: %.2f [#/sec] (mean)\n", float64(*count)/duration)
			fmt.Printf("Time per request: %.2f [ms] (mean)\n", float64(stats.hist.Mean())/1000000)
			fmt.Printf("Fastest request: %.2f [ms]\n", float64(stats.hist.Min())/1000000)
			fmt.Printf("Slowest request: %.2f [ms]\n", float64(stats.hist.Max())/1000000)
			return
		}
	}
}

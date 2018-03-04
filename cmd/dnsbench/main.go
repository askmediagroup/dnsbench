package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/askcom/dnsbench"
	"github.com/codahale/hdrhistogram"
)

func main() {
	var nameserver string
	stats := dnsbench.Stats{}
	stats.Hist = hdrhistogram.New(0, 5e+9, 3)
	stats.IntHist = hdrhistogram.New(0, 5e+9, 3)

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
	results := make(chan *dnsbench.Result, *concur)
	cleanup := make(chan bool, 2)
	names, err := dnsbench.ReadNames(*namesfile)
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
			go dnsbench.LocalResolve(names, t, *qps, results)
		} else {
			go dnsbench.Resolve(nameserver, names, t, *qps, results)
		}
	}

	fmt.Printf("# requests errors min  [ p50  p95  p99  p999] max  qps\n")
	for {
		select {
		case result := <-results:
			stats.Hist.RecordValue(int64(result.Duration))
			stats.IntHist.RecordValue(int64(result.Duration))
			if result.Error {
				stats.Errors++
				stats.IntErrors++
			}
			if stats.Hist.TotalCount() >= int64(*count) {
				cleanup <- true
			}
		case <-timeout:
			dnsbench.DisplayReport(stats.IntHist, stats.IntErrors, float64(*interval)/float64(time.Second))
			stats.IntHist.Reset()
			stats.IntErrors = 0
			timeout = time.After(*interval)
		case <-cleanup:
			end := time.Now().UnixNano()
			duration := float64(end-start) / float64(time.Second)

			fmt.Printf("\nFinished %d requests\n\n", *count)
			fmt.Printf("# latency summary\n")
			dnsbench.DisplayReport(stats.Hist, stats.Errors, duration)
			fmt.Printf("\nConcurrency level: %d\n", *concur)
			fmt.Printf("Time taken for tests: %.2f seconds\n", duration)
			fmt.Printf("Completed Requests: %d\n", int64(*count)-stats.Errors)
			fmt.Printf("Failed Requests: %d\n", stats.Errors)
			fmt.Printf("Requests per second: %.2f [#/sec] (mean)\n", float64(*count)/duration)
			fmt.Printf("Time per request: %.2f [ms] (mean)\n", float64(stats.Hist.Mean())/1000000)
			fmt.Printf("Fastest request: %.2f [ms]\n", float64(stats.Hist.Min())/1000000)
			fmt.Printf("Slowest request: %.2f [ms]\n", float64(stats.Hist.Max())/1000000)
			return
		}
	}
}

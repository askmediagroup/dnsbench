# DNSBench

> A simple DNS server benchmarking tool.

Have you ever needed to troubleshoot a DNS server? Or maybe you want to do some benchmarking before trouble arises. DNSBench is a simple command-line tool that can help you with that.

## Installation

With [Go](https://golang.org/doc/install) and [Dep](https://golang.github.io/dep/docs/introduction.html) installed:

```bash
git clone https://github.com/askcom/dnsbench.git

cd dnsbench

dep init

go build dnsbench.go
```

## Usage and examples

```bash
./dnsbench <nameserver> [flags]
  -concurrency int
        Numer of concurrent requests (default 10)
  -count int
        Number of requests to make (default 1000)
  -interval duration
        Reporting interval (default 5s)
  -local
        Use local resolver
  -names string
        File containing names to request (default "-")
  -qps int
        QPS to attempt for each worker (0 for no limit)
```

Example 1: Benchmark DNS using local resolver:

```bash
$ echo "example.com" | ./dnsbench -local

Reading names from /dev/stdin
Benchmarking...

# requests errors min  [ p50  p95  p99  p999] max  qps

Finished 1000 requests

# latency summary
   1000       0   2.21    [6.31 24.12 36.80 72.35] 72.35  639.33

Concurrency level: 10
Time taken for tests: 1.56 seconds
Completed Requests: 1000
Failed Requests: 0
Requests per second: 639.33 [#/sec] (mean)
Time per request: 8.23 [ms] (mean)
Fastest request: 2.21 [ms]
Slowest request: 72.35 [ms]
```

Example 2: Benchmark a specified nameserver with a file of domains:

```bash
./dnsbench <nameserver ip> -names "domains_to_lookup.txt"

Reading names from domains_to_lookup.txt
Benchmarking...
```

with domains listed on individual lines of **domains_to_lookup.txt**, such as:

```text
example.com
google.com
foobar.com
```

## Status

DNSBench is currently under active development with upcoming improvements targetting:

* Documentation
* Command-line useability
* Decoupling command-line and DNS lookup logic of source
* Testing
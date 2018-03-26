# DNSBench

> A simple DNS server benchmarking tool.

Have you ever needed to troubleshoot a DNS server? Or maybe you want to do some benchmarking before trouble arises. DNSBench is a simple command-line tool that can help you with that.

## Installation

With [Go](https://golang.org/doc/install) installed:

```bash
go get -u github.com/askcom/dnsbench/cmd/dnsbench
```

## Usage and examples

```bash
Execute a latency test.

Usage:
  dnsbench run [command]

Available Commands:
  local       Benchmark the local system resolver configuration.
  remote      Benchmark a remote nameserver.

Flags:
  -c, --concurrency int     Number of concurrent requests. (default 10)
  -n, --count int           Total number requests. (default 1000)
  -h, --help                help for run
  -i, --interval duration   Number of concurrent requests. (default 5s)
  -f, --names string        File containing newline delimited records to lookup. (- for stdin) (default "-")
  -q, --qps int             QPS target for each concurrent worker.

Use "dnsbench run [command] --help" for more information about a command.
```

Example 1: Benchmark DNS using local resolver:

```bash
$ echo "example.com" | dnsbench run local --concurrency 1 --count 10
Reading names from /dev/stdin
Benchmarking...

# requests errors min  [ p50  p95  p99  p999] max  qps

Finished 10 requests

# latency summary
     10       0   0.52    [0.68 8.49 8.49  8.49] 8.49  606.69

Concurrency level: 1
Time taken for tests: 0.02 seconds
Completed Requests: 10
Failed Requests: 0
Requests per second: 606.69 [#/sec] (mean)
Time per request: 1.64 [ms] (mean)
Fastest request: 0.52 [ms]
Slowest request: 8.49 [ms]
```

Example 2: Benchmark a specified nameserver with a file of domains:

```bash
dnsbench run remote 8.8.8.8 --names "domains_to_lookup.txt"
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

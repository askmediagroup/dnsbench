// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"

	dnsbench "github.com/askmediagroup/dnsbench/pkg"
)

func newRunCommand() *cobra.Command {
	opts := runOpts{}
	cmd := &cobra.Command{
		Use:   "run [command]",
		Short: "Execute a DNS benchmark test.",
		RunE: func(_ *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.qps, "qps", "q", 0, "QPS target for the test run. [0 = No limit]")
	cmd.Flags().IntVarP(&opts.count, "count", "c", 100, "Number of queries to attempt. [0 = run until interrupted]")
	cmd.Flags().IntVarP(&opts.maxWorkers, "max-workers", "m", dnsbench.DefaultMaxWorkers, "Maximum number of workers to spawn.")
	cmd.Flags().IntVarP(&opts.workers, "workers", "w", dnsbench.DefaultWorkers, "Initial worker count.")
	cmd.Flags().StringVarP(&opts.resolver, "resolver", "r", "remote", "Resolver mode. [remote,local]")
	cmd.Flags().StringVarP(&opts.nameserver, "nameserver", "n", dnsbench.DefaultNameserver, "Nameserver to benchmark.")
	cmd.Flags().StringVarP(&opts.names, "names", "f", "-", "Read query names from this file.")
	cmd.Flags().DurationVarP(&opts.interval, "interval", "i", 5*time.Second, "Reporting interval.")
	return cmd
}

type runOpts struct {
	nameserver string
	port       int
	qps        int
	count      int
	maxWorkers int
	workers    int
	resolver   string
	interval   time.Duration
	names      string
}

func run(opts runOpts) error {
	if !strings.Contains(opts.nameserver, ":") {
		opts.nameserver += ":53"
	}
	if opts.workers > opts.maxWorkers {
		return fmt.Errorf("initial worker count (%d) must be less than or equal to max worker count (%d)", opts.workers, opts.maxWorkers)
	}
	var resolver dnsbench.Resolver
	if opts.resolver == "remote" {
		resolver = dnsbench.NewRemoteResolver(dnsbench.Nameserver(opts.nameserver))
	} else {
		resolver = dnsbench.NewLocalResolver()
	}
	run := dnsbench.NewRunner(
		resolver,
		dnsbench.Workers(opts.workers),
		dnsbench.MaxWorkers(opts.maxWorkers),
	)
	var file *os.File
	var err error
	if opts.names == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(opts.names)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	namer, err := dnsbench.FileNamer(file)
	if err != nil {
		return err
	}
	res := run.Start(namer, opts.qps, opts.count)
	rep := dnsbench.NewTextReporter(os.Stdout, run.Begin(), opts.interval)
	rep.Heading()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	tick := time.After(time.Second * 5)
L:
	for {
		select {
		case <-sig:
			run.Stop()
			break L
		case r, ok := <-res:
			if !ok {
				break L
			}
			rep.Update(r)
		case <-tick:
			rep.Display()
			tick = time.After(time.Second * 5)
			rep.NewInterval()
		}
	}
	rep.Summary(run)
	return nil
}

package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var concurrency int
var count int
var interval time.Duration
var names string
var qps int

// runCmd verbally extends the root/base command to one or more distinct sub-commands
// e.g., '<base command> run <sub command>' and also propogates shared persistent flags
var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Execute a latency test.",
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent requests.")
	runCmd.PersistentFlags().IntVarP(&count, "count", "n", 1000, "Total number requests.")
	runCmd.PersistentFlags().DurationVarP(&interval, "interval", "i", time.Duration(5)*time.Second, "Number of concurrent requests.")
	runCmd.PersistentFlags().StringVarP(&names, "names", "f", "-", "File containing newline delimited records to lookup. (- for stdin)")
	runCmd.PersistentFlags().IntVarP(&qps, "qps", "q", 0, "QPS target for each concurrent worker.")
}

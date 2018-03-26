package cmd

import (
	"github.com/spf13/cobra"
)

// localCmd represents the local system resolver benchmark command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Benchmark the local system resolver configuration.",
	Run: func(cmd *cobra.Command, args []string) {
		var nameserver string // temp dummy var to accomodate benchmark func signature as-is
		resolveLocally := true
		benchmark(nameserver, &resolveLocally, &concurrency, &count, &interval, &names, &qps)
	},
}

func init() {
	runCmd.AddCommand(localCmd)
}

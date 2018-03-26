package cmd

import (
	"fmt"
	"net"
	"errors"
	"github.com/spf13/cobra"
)

var transport string

// remoteCmd represents the remote nameserver benchmark command
var remoteCmd = &cobra.Command{
	Use:   "remote [nameserver]",
	Short: "Benchmark a remote nameserver.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
		  return errors.New("requires at least one arg [nameserver]")
		}
		if nameserver := net.ParseIP(args[0]); nameserver != nil  {
		  return nil
		}
		return fmt.Errorf("nameserver argument is invalid IP: %s", args[0])
	  },
	Run: func(cmd *cobra.Command, args []string) {
		nameserver := args[0]
		resolveLocally := false
		benchmark(nameserver, &resolveLocally, &concurrency, &count, &interval, &names, &qps)
	},
}

func init() {
	runCmd.AddCommand(remoteCmd)
}

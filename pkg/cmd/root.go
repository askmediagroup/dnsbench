// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root command for the dnsbench command.
func NewRootCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "dnsbench [command]",
		Short: "A simple DNS latency benchmark",
	}
	c.AddCommand(newRunCommand())
	c.AddCommand(newVersionCommand())
	return c
}

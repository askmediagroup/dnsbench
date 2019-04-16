// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	dnsbenchVersion = "unknown"
	gitCommit       = "HEAD"
	buildDate       = "1970-01-01T00:00:00Z"
	goos            = runtime.GOOS
	goarch          = runtime.GOARCH
)

type version struct {
	dnsenchVersion string
	gitCommit      string
	buildDate      string
	goOs           string
	goArch         string
}

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display program version.",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Version: %+v\n", version{
				dnsbenchVersion,
				gitCommit,
				buildDate,
				goos,
				goarch,
			})
		},
	}
}

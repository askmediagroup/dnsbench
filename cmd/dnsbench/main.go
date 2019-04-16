// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/askmediagroup/dnsbench/pkg/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

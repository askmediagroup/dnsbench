// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import (
	"strings"
	"testing"
)

func TestNamer(t *testing.T) {
	names := []string{"github.com", "gitlab.com", "bitbucket.org"}
	expectedNames := []string{"github.com.", "gitlab.com.", "bitbucket.org.", "github.com."}
	reader := strings.NewReader(strings.Join(names, "\n"))

	namer, err := FileNamer(reader)
	if err != nil {
		t.Error(err)
	}

	for i, n := range []Name{namer(), namer(), namer(), namer()} {
		if n.Name != expectedNames[i] {
			t.Errorf("got `%s` for call %d expected `%s`", n.Name, i+1, expectedNames[i])
		}
	}
}

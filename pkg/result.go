// Copyright 2019 AskMediaGroup.
// SPDX-License-Identifier: Apache-2.0

package dnsbench

import "time"

// Result records the diration and error state of a completed DNS request.
type Result struct {
	Duration time.Duration
	Error    error
}

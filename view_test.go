// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatedCursorAndOrigin(t *testing.T) {
	tests := []struct {
		prevOrigin     int
		size           int
		cursor         int
		expectedCursor int
		expectedOrigin int
	}{
		{0, 10, 0, 0, 0},
		{0, 10, 10, 10, 0},
		{0, 10, 11, 10, 1},
		{0, 10, 20, 10, 10},
		{20, 10, 19, 0, 19},
		{20, 10, 25, 5, 20},
	}

	for _, test := range tests {
		cursor, origin := updatedCursorAndOrigin(test.prevOrigin, test.size, test.cursor)
		assert.EqualValues(t, test.expectedCursor, cursor, "Cursor is wrong")
		assert.EqualValues(t, test.expectedOrigin, origin, "Origin in wrong")
	}
}

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

func TestContainsColoredText(t *testing.T) {
	hexColor := func(text string, hexStr string) []cell {
		cells := make([]cell, len(text))
		hex := GetColor(hexStr)
		for i, chr := range text {
			cells[i] = cell{fgColor: hex, chr: chr}
		}
		return cells
	}
	red := "#ff0000"
	green := "#00ff00"
	redStr := func(text string) []cell { return hexColor(text, red) }
	greenStr := func(text string) []cell { return hexColor(text, green) }

	concat := func(lines ...[]cell) []cell {
		var cells []cell
		for _, line := range lines {
			cells = append(cells, line...)
		}
		return cells
	}

	tests := []struct {
		lines      [][]cell
		fgColorStr string
		text       string
		expected   bool
	}{
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "a",
			expected:   true,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: red,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("a"))},
			fgColorStr: green,
			text:       "b",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
		{
			lines:      [][]cell{concat(redStr("hel"), greenStr("lo"), redStr(" World!"))},
			fgColorStr: green,
			text:       "lo",
			expected:   true,
		},
		{
			lines: [][]cell{
				redStr("hel"),
				redStr("lo"),
			},
			fgColorStr: red,
			text:       "hello",
			expected:   false,
		},
	}

	for i, test := range tests {
		v := &View{lines: test.lines}
		assert.Equal(t, test.expected, v.ContainsColoredText(test.fgColorStr, test.text), "Test %d failed", i)
	}
}

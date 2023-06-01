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
	color := func(text string, color Attribute) []cell {
		cells := make([]cell, len(text))
		for i, chr := range text {
			cells[i] = cell{fgColor: color, chr: chr}
		}
		return cells
	}
	red := func(text string) []cell {
		return color(text, ColorRed)
	}
	green := func(text string) []cell {
		return color(text, ColorGreen)
	}

	concat := func(lines ...[]cell) []cell {
		var cells []cell
		for _, line := range lines {
			cells = append(cells, line...)
		}
		return cells
	}

	tests := []struct {
		lines    [][]cell
		color    Attribute
		text     string
		expected bool
	}{
		{
			lines:    [][]cell{concat(red("a"))},
			color:    ColorRed,
			text:     "a",
			expected: true,
		},
		{
			lines:    [][]cell{concat(red("a"))},
			color:    ColorRed,
			text:     "b",
			expected: false,
		},
		{
			lines:    [][]cell{concat(red("a"))},
			color:    ColorGreen,
			text:     "b",
			expected: false,
		},
		{
			lines:    [][]cell{concat(red("hel"), green("lo"), red(" World!"))},
			color:    ColorRed,
			text:     "hello",
			expected: false,
		},
		{
			lines:    [][]cell{concat(red("hel"), green("lo"), red(" World!"))},
			color:    ColorGreen,
			text:     "lo",
			expected: true,
		},
		{
			lines: [][]cell{
				red("hel"),
				red("lo"),
			},
			color:    ColorRed,
			text:     "hello",
			expected: false,
		},
	}

	for _, test := range tests {
		v := &View{lines: test.lines}
		assert.EqualValues(t, test.expected, v.ContainsColoredText(test.color, test.text))
	}
}

// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"github.com/nsf/termbox-go"
)

// A View is a window. It maintains its own internal buffer and cursor
// position.
type View struct {
	name                   string
	x0, y0, x1, y1         int
	ox, oy                 int
	cx, cy                 int
	lines                  [][]rune
	bgColor, fgColor       Attribute
	selBgColor, selFgColor Attribute
	overwrite              bool // overwrite in edit mode

	// If Editable is true, keystrokes will be added to the view's internal
	// buffer at the cursor position.
	Editable bool

	// If Highlight is true, Sel{Bg,Fg}Colors will be used
	// for the line under the cursor position.
	Highlight bool
}

// newView returns a new View object.
func newView(name string, x0, y0, x1, y1 int) *View {
	v := &View{
		name: name,
		x0:   x0,
		y0:   y0,
		x1:   x1,
		y1:   y1,
	}
	return v
}

// Size returns the number of visible columns and rows in the View.
func (v *View) Size() (x, y int) {
	return v.x1 - v.x0 - 1, v.y1 - v.y0 - 1
}

// Name returns the name of the view.
func (v *View) Name() string {
	return v.name
}

// setRune writes a rune at the given point, relative to the view. It
// checks if the position is valid and applies the view's colors, taking
// into account if the cell must be highlighted.
func (v *View) setRune(x, y int, ch rune) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}

	var fgColor, bgColor Attribute
	if v.Highlight && y == v.cy {
		fgColor = v.selFgColor
		bgColor = v.selBgColor
	} else {
		fgColor = v.fgColor
		bgColor = v.bgColor
	}
	termbox.SetCell(v.x0+x+1, v.y0+y+1, ch,
		termbox.Attribute(fgColor), termbox.Attribute(bgColor))
	return nil
}

// SetCursor sets the cursor position of the view at the given point,
// relative to the view. It checks if the position is valid.
func (v *View) SetCursor(x, y int) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}
	v.cx = x
	v.cy = y
	return nil
}

// Cursor returns the cursor position of the view.
func (v *View) Cursor() (x, y int) {
	return v.cx, v.cy
}

// SetOrigin sets the origin position of the view's internal buffer,
// so the buffer starts to be printed from this point, which means that
// it is linked with the origin point of view. It can be used to
// implement Horizontal and Vertical scrolling with just incrementing
// or decrementing ox and oy.
func (v *View) SetOrigin(x, y int) error {
	if x < 0 || y < 0 {
		return errors.New("invalid point")
	}
	v.ox = x
	v.oy = y
	return nil
}

// Origin returns the origin position of the view.
func (v *View) Origin() (x, y int) {
	return v.ox, v.oy
}

// Write appends a byte slice into the view's internal buffer. Because
// View implements the io.Writer interface, it can be passed as parameter
// of functions like fmt.Fprintf, fmt.Fprintln, io.Copy, etc. Clear must
// be called to clear the view's buffer.
func (v *View) Write(p []byte) (n int, err error) {
	r := bytes.NewReader(p)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := bytes.Runes(s.Bytes())
		v.lines = append(v.lines, line)
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return len(p), nil
}

// draw re-draws the view's contents.
func (v *View) draw() error {
	maxX, maxY := v.Size()
	y := 0
	for i, line := range v.lines {
		if i < v.oy {
			continue
		}
		x := 0
		for j, ch := range line {
			if j < v.ox {
				continue
			}
			if x >= 0 && x < maxX && y >= 0 && y < maxY {
				if err := v.setRune(x, y, ch); err != nil {
					return err
				}
			}
			x++
		}
		y++
	}
	return nil
}

// Clear empties the view's internal buffer.
func (v *View) Clear() {
	v.lines = nil
	v.clearRunes()
}

// clearRunes erases all the cells in the view.
func (v *View) clearRunes() {
	maxX, maxY := v.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			termbox.SetCell(v.x0+x+1, v.y0+y+1, ' ',
				termbox.Attribute(v.fgColor), termbox.Attribute(v.bgColor))
		}
	}
}

// writeRune writes a rune into the view's internal buffer, at the
// position corresponding to the point (x, y). The length of the internal
// buffer is increased if the point is out of bounds. Overwrite mode is
// governed by the value of View.overwrite.
func (v *View) writeRune(x, y int, ch rune) error {
	x = v.ox + x
	y = v.oy + y

	if x < 0 || y < 0 {
		return errors.New("invalid point")
	}

	if y >= len(v.lines) {
		if y >= cap(v.lines) {
			s := make([][]rune, y+1, (y+1)*2)
			copy(s, v.lines)
			v.lines = s
		} else {
			v.lines = v.lines[:y+1]
		}
	}
	if v.lines[y] == nil {
		v.lines[y] = make([]rune, x+1, (x+1)*2)
	} else if x >= len(v.lines[y]) {
		if x >= cap(v.lines[y]) {
			s := make([]rune, x+1, (x+1)*2)
			copy(s, v.lines[y])
			v.lines[y] = s
		} else {
			v.lines[y] = v.lines[y][:x+1]
		}
	}
	if !v.overwrite {
		v.lines[y] = append(v.lines[y], ' ')
		copy(v.lines[y][x+1:], v.lines[y][x:])
	}
	v.lines[y][x] = ch
	return nil
}

// deleteRune removes a rune from the view's internal buffer, at the
// position corresponding to the point (x, y).
func (v *View) deleteRune(x, y int) error {
	x = v.ox + x
	y = v.oy + y

	if x < 0 || y < 0 || y >= len(v.lines) || v.lines[y] == nil || x >= len(v.lines[y]) {
		return errors.New("invalid point")
	}
	copy(v.lines[y][x:], v.lines[y][x+1:])
	v.lines[y][len(v.lines[y])-1] = ' '
	return nil
}

// addLine adds a line into the view's internal buffer at the position
// corresponding to the point (x, y).
func (v *View) addLine(y int) error {
	y = v.oy + y

	if y < 0 || y >= len(v.lines) {
		return errors.New("invalid point")
	}
	v.lines = append(v.lines, nil)
	copy(v.lines[y+1:], v.lines[y:])
	v.lines[y] = nil
	return nil
}

// Line returns a string with the line of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Line(y int) (string, error) {
	y = v.oy + y

	if y < 0 || y >= len(v.lines) {
		return "", errors.New("invalid point")
	}
	return string(v.lines[y]), nil
}

// Word returns a string with the word of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Word(x, y int) (string, error) {
	x = v.ox + x
	y = v.oy + y

	if y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
		return "", errors.New("invalid point")
	}
	l := string(v.lines[y])
	nl := strings.LastIndexFunc(l[:x], indexFunc)
	if nl == -1 {
		nl = 0
	} else {
		nl = nl + 1
	}
	nr := strings.IndexFunc(l[x:], indexFunc)
	if nr == -1 {
		nr = len(l)
	} else {
		nr = nr + x
	}
	return string(l[nl:nr]), nil
}

// indexFunc allows to split lines by words taking into account spaces
// and 0
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

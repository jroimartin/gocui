// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

// A View is a window. It maintains its own internal buffer and cursor
// position.
type View struct {
	name                   string
	x0, y0, x1, y1         int
	ox, oy                 int
	cx, cy, realx, realy   int
	lines                  [][]rune
	drawpos2linepos        []map[int]int
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
func (v *View) getRuneLen(ch rune) int {
	if utf8.RuneLen(ch) > 1 {
		return 2
	} else {
		return 1
	}
}
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
	le := v.getRuneLen(ch)
	for i:=1; i < le; i++ {
		termbox.SetCell(v.x0+x+1+i, v.y0+y+1, 0,
			termbox.Attribute(fgColor), termbox.Attribute(bgColor))
	}
	return nil
}

// SetCursor sets the cursor position of the view at the given point,
// relative to the view. It checks if the position is valid.
func (v *View) SetCursor(x, y int) error {
	oldrealx := v.realx
	maxX, maxY := v.Size()
	destRealX := x + v.ox
	destRealY := y + v.oy
	totalL := v.getLineTail(destRealY)
	if v.getIndexCheck(destRealX, destRealY) == -1 {
		if destRealY < 0 {
			return v.SetCursor(x, y+1)
		} else if destRealY >= len(v.lines) {
			return v.SetCursor(x, y-1)
		}
		if destRealX < totalL {
			if oldrealx < destRealX {
				return v.SetCursor(x+1, y)
			} else if destRealX > 0 {
				return v.SetCursor(x-1, y)
			}
		} else if destRealX > totalL {
			return v.SetCursor(x-1, y)
		}
	}
	viewX := destRealX - v.ox
	if viewX < 0 {
		if v.ox > 0 {
			v.ox--
			return v.SetCursor(x+1, y)
		} else {
			viewX = 0
		}
	} else if viewX >= maxX {
		v.ox++
		return v.SetCursor(x-1, y)
	}
	viewY := destRealY - v.oy
	if viewY < 0 {
		if v.oy > 0 {
			v.oy--
			return v.SetCursor(x, y+1)
		} else {
			viewY = 0
		}
	} else if viewY >= maxY {
		v.oy++
		return v.SetCursor(x, y-1)
	}
	if v.ox < 0 {
		viewX = v.cx
		v.ox = 0
	}
	if v.oy < 0 {
		viewY = v.cy
		v.oy = 0
	}
	v.cx = viewX
	v.cy = viewY
	v.realx = viewX + v.ox
	v.realy = viewY + v.oy
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
func (v *View) resetDrawPos(y int) map[int]int {
	line := v.lines[y]
	drawp2linep := make(map[int]int)
	base := 0
	for i, ch := range line {
		drawp2linep[base] = i
		base += v.getRuneLen(ch)
	}
	v.drawpos2linepos[y] = drawp2linep
	return drawp2linep
}

func (v *View) Write(p []byte) (n int, err error) {
	r := bytes.NewReader(p)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := bytes.Runes(s.Bytes())
		v.lines = append(v.lines, line)
		v.drawpos2linepos = append(v.drawpos2linepos, make(map[int]int))
		v.resetDrawPos(len(v.lines) - 1)
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
		for _, ch := range line {
			if x < v.ox {
				x += v.getRuneLen(ch)
				continue
			}
			offset := x - v.ox
			if offset >= 0 && offset < maxX && y >= 0 && y < maxY {
				if err := v.setRune(offset, y, ch); err != nil {
					return err
				}
				x += v.getRuneLen(ch)
			}
		}
		y++
	}
	return nil
}

// Clear empties the view's internal buffer.
func (v *View) Clear() {
	v.lines = nil
	v.drawpos2linepos = nil
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

	orilen := len(v.lines[y])
	if y >= len(v.lines) {
		if y >= cap(v.lines) {
			s := make([][]rune, y+1, (y+1)*2)
			copy(s, v.lines)
			v.lines = s

			t := make([]map[int]int, y+1, (y+1)*2)
			copy(t, v.drawpos2linepos)
			v.drawpos2linepos = t
		} else {
			v.lines = v.lines[:y+1]
			v.drawpos2linepos = v.drawpos2linepos[:y+1]
		}
	}
	currpos := x
	if v.lines[y] == nil {
		v.lines[y] = make([]rune, x+1, (x+1)*2)
		v.drawpos2linepos[y] = make(map[int]int)
	} else {
		currpos = v.getIndex(x, y)
	}
	if currpos >= len(v.lines[y]) {
		if currpos >= cap(v.lines[y]) {
			s := make([]rune, currpos+1, (currpos+1)*2)
			copy(s, v.lines[y])
			v.lines[y] = s
		} else {
			v.lines[y] = v.lines[y][:currpos+1]
		}
	}
	if !v.overwrite {
		if currpos < orilen {
			v.lines[y] = append(v.lines[y], ' ')
		}
		copy(v.lines[y][currpos+1:], v.lines[y][currpos:])
	}
	v.lines[y][currpos] = ch
	v.resetDrawPos(y)
	v.SetCursor(v.cx+v.getRuneLen(ch), v.cy)
	return nil
}

// deleteRune removes a rune from the view's internal buffer, at the
// position corresponding to the point (x, y).
func (v *View) deleteRune(x, y int) int {
	if y < 0 {
		return 0
	}
	realx := x + v.ox
	realy := y + v.oy
	if realx < 0 && realy > 0 {
		l := v.getLineWidth(realy)
		realy--
		g := v.getLineTail(realy)
		v.lines[realy] = append(v.lines[realy], v.lines[realy+1]...)
		v.resetDrawPos(realy)
		if g == -1 {
			g = 0
		}
		v.SetCursor(g-v.ox, y-1)
		le := len(v.lines)
		if le > realy+2 {
			copy(v.lines[realy+1:], v.lines[realy+2:])
			v.lines = v.lines[:le-1]
			copy(v.drawpos2linepos[realy+1:], v.drawpos2linepos[realy+2:])
			v.drawpos2linepos = v.drawpos2linepos[:le-1]
		} else {
			v.lines = v.lines[:realy+1]
			v.drawpos2linepos = v.drawpos2linepos[:realy+1]
		}
		return l
	}
	x = realx
	y = realy

	currpos := v.getIndexCheck(x, y)
	if currpos == -1 && x > 0 {
		return v.deleteRune(x-v.ox-1, y-v.oy)
	}
	if x < 0 || y < 0 || y >= len(v.lines) || v.lines[y] == nil || currpos >= len(v.lines[y]) {
		return 0
	}
	l := v.getRuneLen(v.lines[y][currpos])
	copy(v.lines[y][currpos:], v.lines[y][currpos+1:])
	v.lines[y] = v.lines[y][:len(v.lines[y])-1]
	v.resetDrawPos(y)
	v.SetCursor(x-v.ox, y-v.oy)
	return l
}

func (v *View) RemoveTextToLineEnd(x, y int) string {
	realx, realy := x+v.ox, y+v.oy
	if realy < 0 || realy >= len(v.lines) {
		return ""
	}
	tail := v.GetTextToLineEnd(x, y)
	if realx <= 0 {
		v.lines[realy] = v.lines[realy]
	} else if realx >= v.getLineTail(realy) {
		v.lines[realy] = []rune{}
	} else {
		v.lines[realy] = v.lines[realy][:(realx - 1)]
	}
	v.resetDrawPos(realy)
	return tail
}

func (v *View) GetTextToLineBegin(x, y int) string {
	realx, realy := x+v.ox, y+v.oy
	if realy < 0 || realy >= len(v.lines) {
		return ""
	}
	line := v.lines[realy]
	pos, bHave := v.drawpos2linepos[realy][realx]
	if bHave {
		return string(line[:pos])
	} else {
		return string(line)
	}
}

func (v *View) GetTextToLineEnd(x, y int) string {
	realx, realy := x+v.ox, y+v.oy
	if realy < 0 || realy >= len(v.lines) {
		return ""
	}
	line := v.lines[realy]
	pos, bHave := v.drawpos2linepos[realy][realx]
	if bHave {
		return string(line[pos:])
	} else {
		return ""
	}
}

// addLine adds a line into the view's internal buffer at the position
// corresponding to the point (x, y).
func (v *View) AddLine(x, y int) error {
	head := v.GetTextToLineBegin(x, y)
	tail := v.RemoveTextToLineEnd(x, y)
	y = v.oy + y

	v.lines[y] = []rune(tail)
	v.lines = append(v.lines, nil)
	copy(v.lines[y+1:], v.lines[y:])
	v.lines[y] = []rune(head)
	v.drawpos2linepos = append(v.drawpos2linepos, nil)
	copy(v.drawpos2linepos[y+1:], v.drawpos2linepos[y:])
	v.resetDrawPos(y)
	v.resetDrawPos(y + 1)
	v.ox = 0
	v.SetCursor(0, y-v.oy+1)
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

func (v *View) getIndex(x, y int) int {
	currpos := 0
	if y < 0 || y >= len(v.drawpos2linepos) {
		return currpos
	}
	linepos, bhave := v.drawpos2linepos[y][x]
	if bhave {
		currpos = linepos
	} else {
		currpos = len(v.lines[y])
	}
	return currpos
}

func (v *View) getLineTail(y int) int {
	if y < 0 || y >= len(v.drawpos2linepos) {
		return 0
	}
	l := -1
	var c rune
	for pos, i := range v.drawpos2linepos[y] {
		if l < pos {
			l = pos
			c = v.lines[y][i]
		}
	}
	return l + v.getRuneLen(c)
}

func (v *View) getLineWidth(y int) int {
	if y < 0 || y >= len(v.drawpos2linepos) {
		return 0
	}
	l := -1
	for pos, _ := range v.drawpos2linepos[y] {
		if l < pos {
			l = pos
		}
	}
	return l
}

func (v *View) getIndexCheck(x, y int) int {
	currpos := -1
	if y < 0 || y >= len(v.drawpos2linepos) {
		return currpos
	}
	linepos, bhave := v.drawpos2linepos[y][x]
	if bhave {
		currpos = linepos
	} else {
		currpos = -1
	}
	return currpos
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
	currpos := v.getIndex(x, y)
	nr := strings.IndexFunc(l[currpos:], indexFunc)
	if nr == -1 {
		nr = len(l)
	} else {
		nr = nr + currpos
	}
	return string(l[nl:nr]), nil
}

// indexFunc allows to split lines by words taking into account spaces
// and 0
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

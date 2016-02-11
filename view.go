// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

// A View is a window. It maintains its own internal buffer and cursor
// position.
type View struct {
	name           string
	x0, y0, x1, y1 int
	ox, oy         int
	cx, cy         int
	lines          [][]rune
	readOffset     int
	readCache      string

	tainted   bool       // marks if the viewBuffer must be updated
	viewLines []viewLine // internal representation of the view's buffer

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the View.
	BgColor, FgColor Attribute

	// SelBgColor and SelFgColor are used to configure the background and
	// foreground colors of the selected line, when it is highlighted.
	SelBgColor, SelFgColor Attribute

	// If Editable is true, keystrokes will be added to the view's internal
	// buffer at the cursor position.
	Editable bool

	// Overwrite enables or disables the overwrite mode of the view.
	Overwrite bool

	// If Highlight is true, Sel{Bg,Fg}Colors will be used
	// for the line under the cursor position.
	Highlight bool

	// If Frame is true, a border will be drawn around the view.
	Frame bool

	// If Wrap is true, the content that is written to this View is
	// automatically wrapped when it is longer than its width. If true the
	// view's x-origin will be ignored.
	Wrap bool

	// If Autoscroll is true, the View will automatically scroll down when the
	// text overflows. If true the view's y-origin will be ignored.
	Autoscroll bool

	// If Frame is true, Title allows to configure a title for the view.
	Title string
}

type viewLine struct {
	linesX, linesY int // coordinates relative to v.lines
	line           []rune
}

// newView returns a new View object.
func newView(name string, x0, y0, x1, y1 int) *View {
	v := &View{
		name:    name,
		x0:      x0,
		y0:      y0,
		x1:      x1,
		y1:      y1,
		Frame:   true,
		tainted: true,
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
// checks if the position is valid
func (v *View) setRune(x, y int, ch rune, fgColor, bgColor Attribute) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
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
	v.tainted = true

	for _, ch := range bytes.Runes(p) {
		switch ch {
		case '\n':
			v.lines = append(v.lines, nil)
		case '\r':
			nl := len(v.lines)
			if nl > 0 {
				v.lines[nl-1] = nil
			} else {
				v.lines = make([][]rune, 1)
			}
		default:
			nl := len(v.lines)
			if nl > 0 {
				v.lines[nl-1] = append(v.lines[nl-1], ch)
			} else {
				v.lines = append(v.lines, []rune{ch})
			}
		}
	}
	return len(p), nil
}

// Read reads data into p. It returns the number of bytes read into p.
// At EOF, err will be io.EOF. Calling Read() after Rewind() makes the
// cache to be refreshed with the contents of the view.
func (v *View) Read(p []byte) (n int, err error) {
	if v.readOffset == 0 {
		v.readCache = v.Buffer()
	}
	if v.readOffset < len(v.readCache) {
		n = copy(p, v.readCache[v.readOffset:])
		v.readOffset += n
	} else {
		err = io.EOF
	}
	return
}

// Rewind sets the offset for the next Read to 0, which also refresh the
// read cache.
func (v *View) Rewind() {
	v.readOffset = 0
}

type escapeState int
type interpreterReturn int

const (
	NONE escapeState = iota
	ESCAPE
	CSI
	PARAMS
)

const (
	ERROR interpreterReturn = iota
	IS_ESCAPE
	IS_NOT_ESCAPE
)

type escapeInterpreter struct {
	state                          escapeState
	curch                          rune
	csi_param                      []string
	curFgColor, curBgColor         Attribute
	defaultFgColor, defaultBgColor Attribute
}

// In case of error, this will output the non-parsed runes as a string
func (ei *escapeInterpreter) runes() []rune {
	switch ei.state {
	case NONE:
		return []rune{0x1b}
	case ESCAPE:
		return []rune{0x1b, ei.curch}
	case CSI:
		return []rune{0x1b, '[', ei.curch}
	case PARAMS:
		ret := []rune{0x1b, '['}
		for _, s := range ei.csi_param {
			ret = append(ret, []rune(s)...)
			ret = append(ret, ';')
		}
		return append(ret, ei.curch)
	}
	return []rune{}
}

//returns an escapeInterpreter that will be able to parse terminal escape sequences
func NewEscapeInterpreter(fg, bg Attribute) *escapeInterpreter {
	ei := &escapeInterpreter{}
	ei.defaultFgColor = fg
	ei.defaultBgColor = bg
	ei.state = NONE
	ei.curFgColor = ei.defaultFgColor
	ei.curBgColor = ei.defaultBgColor
	ei.csi_param = make([]string, 0)
	return ei
}

func (ei *escapeInterpreter) reset() {
	ei.state = NONE
	ei.curFgColor = ei.defaultFgColor
	ei.curBgColor = ei.defaultBgColor
	ei.csi_param = make([]string, 0)
}

var ErrorNotCSI = errors.New("Not a CSI escape sequence")
var ErrorCSINotANumber = errors.New("CSI escape sequence was expecting a number or a ;")
var ErrorCSIParseError = errors.New("CSI escape sequence parsing error")
var ErrorCSITooLong = errors.New("CSI escape sequence is too long")

//returns an attribute given a terminfo coloring
func paramToColor(p int) Attribute {
	switch p {
	case 0:
		return ColorBlack
	case 1:
		return ColorRed
	case 2:
		return ColorGreen
	case 3:
		return ColorYellow
	case 4:
		return ColorBlue
	case 5:
		return ColorMagenta
	case 6:
		return ColorCyan
	case 7:
		return ColorWhite
	}
	return ColorDefault
}

//Parses a rune, returns ERROR, IS_ESCAPE or IS_NOT_ESCAPE
//In case of ERROR, an error is returned to specify the type of error
//In case the return is IS_ESCAPE, it means that the rune is part of an
// scape sequence, and as such should not be printed verbatim
//In case the return is IS_NOT_ESCAPE, it means it's not an escape sequence
func (ei *escapeInterpreter) parseOne(ch rune) (interpreterReturn, error) {
	/* A couple of sanity checks, too make sure we're not parsing
	something totally bogus */
	if len(ei.csi_param) > 20 {
		return ERROR, ErrorCSITooLong
	}
	if len(ei.csi_param) > 0 && len(ei.csi_param[len(ei.csi_param)-1]) > 255 {
		return ERROR, ErrorCSITooLong
	}
	ei.curch = ch
	switch ei.state {
	case NONE:
		if ch == 0x1b {
			ei.state = ESCAPE
			return IS_ESCAPE, nil
		}
		return IS_NOT_ESCAPE, nil
	case ESCAPE:
		if ch == '[' {
			ei.state = CSI
			return IS_ESCAPE, nil
		}
		return ERROR, ErrorNotCSI
	case CSI:
		if ch >= '0' && ch <= '9' {
			ei.state = PARAMS
			ei.csi_param = append(ei.csi_param, string(ch))
			return IS_ESCAPE, nil
		}
		return ERROR, ErrorCSINotANumber
	case PARAMS:
		switch {
		case ch >= '0' && ch <= '9':
			ei.csi_param[len(ei.csi_param)-1] += string(ch)
			return IS_ESCAPE, nil
		case ch == ';':
			ei.csi_param = append(ei.csi_param, "")
			return IS_ESCAPE, nil
		case ch == 'm':
			if len(ei.csi_param) < 1 {
				return ERROR, ErrorCSIParseError
			}
			for _, param := range ei.csi_param {
				p, err := strconv.Atoi(param)
				if err != nil {
					return ERROR, ErrorCSIParseError
				}
				switch {
				case p >= 30 && p <= 37:
					ei.curFgColor = paramToColor(p - 30)
				case p >= 40 && p <= 47:
					ei.curBgColor = paramToColor(p - 40)
				case p == 1:
					ei.curFgColor |= AttrBold
				case p == 4:
					ei.curFgColor |= AttrUnderline
				case p == 7:
					ei.curFgColor |= AttrReverse
				case p == 0 || p == 39:
					ei.curFgColor = ei.defaultFgColor
					ei.curBgColor = ei.defaultBgColor
				}
			}
			ei.state = NONE
			ei.csi_param = make([]string, 0)
			return IS_ESCAPE, nil
		}
	}
	return IS_NOT_ESCAPE, nil
}

//prints the runes in s, it parseEscape is true, we'll try to interpret
//terminal escape sequences
func printLine(v *View, s []rune, maxX, y int, parseEscape bool) error {
	x := 0
	ei := NewEscapeInterpreter(v.FgColor, v.BgColor)
	for j, ch := range s {
		var fgColor, bgColor Attribute
		if j < v.ox {
			continue
		}
		if x >= maxX {
			break
		}

		if parseEscape {
			ret, _ := ei.parseOne(ch)
			switch ret {
			case ERROR:
				printLine(v, ei.runes(), maxX, y, false)
				ei.reset()
			case IS_ESCAPE:
				continue
			}
		}

		if v.Highlight && y == v.cy {
			fgColor = v.SelFgColor
			bgColor = v.SelBgColor
		} else {
			fgColor = ei.curFgColor
			bgColor = ei.curBgColor
		}
		if err := v.setRune(x, y, ch, fgColor, bgColor); err != nil {
			return err
		}
		x++
	}
	return nil
}

// draw re-draws the view's contents.
func (v *View) draw() error {
	maxX, maxY := v.Size()

	if v.Wrap {
		if maxX == 0 {
			return errors.New("X size of the view cannot be 0")
		}
		v.ox = 0
	}
	if v.tainted {
		v.viewLines = nil
		for i, line := range v.lines {
			if v.Wrap {
				if len(line) <= maxX {
					vline := viewLine{linesX: 0, linesY: i, line: line}
					v.viewLines = append(v.viewLines, vline)
					continue
				} else {
					vline := viewLine{linesX: 0, linesY: i, line: line[:maxX]}
					v.viewLines = append(v.viewLines, vline)
				}
				// Append remaining lines
				for n := maxX; n < len(line); n += maxX {
					if len(line[n:]) <= maxX {
						vline := viewLine{linesX: n, linesY: i, line: line[n:]}
						v.viewLines = append(v.viewLines, vline)
					} else {
						vline := viewLine{linesX: n, linesY: i, line: line[n : n+maxX]}
						v.viewLines = append(v.viewLines, vline)
					}
				}
			} else {
				vline := viewLine{linesX: 0, linesY: i, line: line}
				v.viewLines = append(v.viewLines, vline)
			}
		}
		v.tainted = false
	}

	if v.Autoscroll && len(v.viewLines) > maxY {
		v.oy = len(v.viewLines) - maxY
	}
	y := 0

	for i, vline := range v.viewLines {
		if i < v.oy {
			continue
		}
		if y >= maxY {
			break
		}
		printLine(v, vline.line, maxX, y, true)
		y++
	}
	return nil
}

// realPosition returns the position in the internal buffer corresponding to the
// point (x, y) of the view.
func (v *View) realPosition(vx, vy int) (x, y int, err error) {
	vx = v.ox + vx
	vy = v.oy + vy

	if vx < 0 || vy < 0 {
		return 0, 0, errors.New("invalid point")
	}

	if len(v.viewLines) == 0 {
		return vx, vy, nil
	}

	if vy < len(v.viewLines) {
		vline := v.viewLines[vy]
		x = vline.linesX + vx
		y = vline.linesY
	} else {
		vline := v.viewLines[len(v.viewLines)-1]
		x = vx
		y = vline.linesY + vy - len(v.viewLines) + 1
	}

	return x, y, nil
}

// Clear empties the view's internal buffer.
func (v *View) Clear() {
	v.tainted = true

	v.lines = nil
	v.clearRunes()
}

// clearRunes erases all the cells in the view.
func (v *View) clearRunes() {
	maxX, maxY := v.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			termbox.SetCell(v.x0+x+1, v.y0+y+1, ' ',
				termbox.Attribute(v.FgColor), termbox.Attribute(v.BgColor))
		}
	}
}

// writeRune writes a rune into the view's internal buffer, at the
// position corresponding to the point (x, y). The length of the internal
// buffer is increased if the point is out of bounds. Overwrite mode is
// governed by the value of View.overwrite.
func (v *View) writeRune(x, y int, ch rune) error {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return err
	}

	if x < 0 || y < 0 {
		return errors.New("invalid point")
	}

	if y >= len(v.lines) {
		s := make([][]rune, y-len(v.lines)+1)
		v.lines = append(v.lines, s...)
	}

	olen := len(v.lines[y])
	if x >= len(v.lines[y]) {
		s := make([]rune, x-len(v.lines[y])+1)
		v.lines[y] = append(v.lines[y], s...)
	}

	if !v.Overwrite && x < olen {
		v.lines[y] = append(v.lines[y], '\x00')
		copy(v.lines[y][x+1:], v.lines[y][x:])
	}
	v.lines[y][x] = ch
	return nil
}

// deleteRune removes a rune from the view's internal buffer, at the
// position corresponding to the point (x, y).
func (v *View) deleteRune(x, y int) error {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return err
	}

	if x < 0 || y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
		return errors.New("invalid point")
	}
	v.lines[y] = append(v.lines[y][:x], v.lines[y][x+1:]...)
	return nil
}

// mergeLines merges the lines "y" and "y+1" if possible.
func (v *View) mergeLines(y int) error {
	v.tainted = true

	_, y, err := v.realPosition(0, y)
	if err != nil {
		return err
	}

	if y < 0 || y >= len(v.lines) {
		return errors.New("invalid point")
	}

	if y < len(v.lines)-1 { // otherwise we don't need to merge anything
		v.lines[y] = append(v.lines[y], v.lines[y+1]...)
		v.lines = append(v.lines[:y+1], v.lines[y+2:]...)
	}
	return nil
}

// breakLine breaks a line of the internal buffer at the position corresponding
// to the point (x, y).
func (v *View) breakLine(x, y int) error {
	v.tainted = true

	x, y, err := v.realPosition(x, y)
	if err != nil {
		return err
	}

	if y < 0 || y >= len(v.lines) {
		return errors.New("invalid point")
	}

	var left, right []rune
	if x < len(v.lines[y]) { // break line
		left = make([]rune, len(v.lines[y][:x]))
		copy(left, v.lines[y][:x])
		right = make([]rune, len(v.lines[y][x:]))
		copy(right, v.lines[y][x:])
	} else { // new empty line
		left = v.lines[y]
	}

	lines := make([][]rune, len(v.lines)+1)
	lines[y] = left
	lines[y+1] = right
	copy(lines, v.lines[:y])
	copy(lines[y+2:], v.lines[y+1:])
	v.lines = lines
	return nil
}

// Buffer returns a string with the contents of the view's internal
// buffer.
func (v *View) Buffer() string {
	str := ""
	for _, l := range v.lines {
		str += string(l) + "\n"
	}
	return strings.Replace(str, "\x00", " ", -1)
}

// ViewBuffer returns a string with the contents of the view's buffer that is
// shown to the user.
func (v *View) ViewBuffer() string {
	str := ""
	for _, l := range v.viewLines {
		str += string(l.line) + "\n"
	}
	return strings.Replace(str, "\x00", " ", -1)
}

// Line returns a string with the line of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Line(y int) (string, error) {
	_, y, err := v.realPosition(0, y)
	if err != nil {
		return "", err
	}

	if y < 0 || y >= len(v.lines) {
		return "", errors.New("invalid point")
	}
	return string(v.lines[y]), nil
}

// Word returns a string with the word of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Word(x, y int) (string, error) {
	x, y, err := v.realPosition(x, y)
	if err != nil {
		return "", err
	}

	if x < 0 || y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
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
// and 0.
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

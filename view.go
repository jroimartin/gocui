// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// Constants for overlapping edges
const (
	TOP    = 1 // view is overlapping at top edge
	BOTTOM = 2 // view is overlapping at bottom edge
	LEFT   = 4 // view is overlapping at left edge
	RIGHT  = 8 // view is overlapping at right edge
)

var (
	// ErrInvalidPoint is returned when client passed invalid coordinates of a cell.
	// Most likely client has passed negative coordinates of a cell.
	ErrInvalidPoint = errors.New("invalid point")
)

// A View is a window. It maintains its own internal buffer and cursor
// position.
type View struct {
	name           string
	x0, y0, x1, y1 int      // left top right bottom
	ox, oy         int      // view offsets
	cx, cy         int      // cursor position
	rx, ry         int      // Read() offsets
	wx, wy         int      // Write() offsets
	lines          [][]cell // All the data
	outMode        OutputMode

	// readBuffer is used for storing unread bytes
	readBuffer []byte

	// tained is true if the viewLines must be updated
	tainted bool

	// contentCache is the content the frame
	// if a redraw is request with tainted is false this will be used to draw the frame
	contentCache []cellCache

	// writeMutex protects locks the write process
	writeMutex sync.Mutex

	// ei is used to decode ESC sequences on Write
	ei *escapeInterpreter

	// Visible specifies whether the view is visible.
	Visible bool

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the View.
	BgColor, FgColor Attribute

	// SelBgColor and SelFgColor are used to configure the background and
	// foreground colors of the selected line, when it is highlighted.
	SelBgColor, SelFgColor Attribute

	// If Editable is true, keystrokes will be added to the view's internal
	// buffer at the cursor position.
	Editable bool

	// Editor allows to define the editor that manages the editing mode,
	// including keybindings or cursor behaviour. DefaultEditor is used by
	// default.
	Editor Editor

	// Overwrite enables or disables the overwrite mode of the view.
	Overwrite bool

	// If Highlight is true, Sel{Bg,Fg}Colors will be used
	// for the line under the cursor position.
	Highlight bool

	// If Frame is true, a border will be drawn around the view.
	Frame bool

	// FrameColor allow to configure the color of the Frame when it is not highlighted.
	FrameColor Attribute

	// FrameRunes allows to define custom runes for the frame edges.
	// The rune slice can be defined with 3 different lengths.
	// If slice doesn't match these lengths, default runes will be used instead of missing one.
	//
	// 2 runes with only horizontal and vertical edges.
	//  []rune{'─', '│'}
	//  []rune{'═','║'}
	// 6 runes with horizontal, vertical edges and top-left, top-right, bottom-left, bottom-right cornes.
	//  []rune{'─', '│', '┌', '┐', '└', '┘'}
	//  []rune{'═','║','╔','╗','╚','╝'}
	// 11 runes which can be used with `gocui.Gui.SupportOverlaps` property.
	//  []rune{'─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼'}
	//  []rune{'═','║','╔','╗','╚','╝','╠','╣','╦','╩','╬'}
	FrameRunes []rune

	// If Wrap is true, the content that is written to this View is
	// automatically wrapped when it is longer than its width. If true the
	// view's x-origin will be ignored.
	Wrap bool

	// If Autoscroll is true, the View will automatically scroll down when the
	// text overflows. If true the view's y-origin will be ignored.
	Autoscroll bool

	// If Frame is true, Title allows to configure a title for the view.
	Title string

	// TitleColor allow to configure the color of title and subtitle for the view.
	TitleColor Attribute

	// If Frame is true, Subtitle allows to configure a subtitle for the view.
	Subtitle string

	// If Mask is true, the View will display the mask instead of the real
	// content
	Mask rune

	// Overlaps describes which edges are overlapping with another view's edges
	Overlaps byte

	// If HasLoader is true, the message will be appended with a spinning loader animation
	HasLoader bool

	// KeybindOnEdit should be set to true when you want to execute keybindings even when the view is editable
	// (this is usually not the case)
	KeybindOnEdit bool

	// gui contains the view it's gui
	gui *Gui
}

type cell struct {
	chr              rune
	bgColor, fgColor Attribute
}

type cellCache struct {
	chr              rune
	bgColor, fgColor Attribute
	x, y             int
}

type lineType []cell

// String returns a string from a given cell slice.
func (l lineType) String() string {
	str := ""
	for _, c := range l {
		str += string(c.chr)
	}
	return str
}

// newView returns a new View object.
func (g *Gui) newView(name string, x0, y0, x1, y1 int, mode OutputMode) *View {
	v := &View{
		name:    name,
		x0:      x0,
		y0:      y0,
		x1:      x1,
		y1:      y1,
		Visible: true,
		Frame:   true,
		Editor:  DefaultEditor,
		tainted: true,
		outMode: mode,
		ei:      newEscapeInterpreter(mode),
		gui:     g,
	}

	v.FgColor, v.BgColor = ColorDefault, ColorDefault
	v.SelFgColor, v.SelBgColor = ColorDefault, ColorDefault
	v.TitleColor, v.FrameColor = ColorDefault, ColorDefault
	return v
}

// Dimensions returns the dimensions of the View
func (v *View) Dimensions() (int, int, int, int) {
	return v.x0, v.y0, v.x1, v.y1
}

// Size returns the number of visible columns and rows in the View.
func (v *View) Size() (x, y int) {
	return v.x1 - v.x0 - 1, v.y1 - v.y0 - 1
}

// Name returns the name of the view.
func (v *View) Name() string {
	return v.name
}

// setRune sets a rune at the given point relative to the view. It applies the
// specified colors, taking into account if the cell must be highlighted. Also,
// it checks if the position is valid.
func (v *View) setRune(x, y int, ch rune, fgColor, bgColor Attribute) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return ErrInvalidPoint
	}

	if v.Mask != 0 {
		fgColor = v.FgColor
		bgColor = v.BgColor
		ch = v.Mask
	} else if v.Highlight && y == v.cy {
		fgColor = v.SelFgColor | AttrBold
		bgColor = v.SelBgColor | AttrBold
	}

	// Don't display NUL characters
	if ch == 0 {
		ch = ' '
	}

	tcellSetCell(v.x0+x+1, v.y0+y+1, ch, fgColor, bgColor, v.outMode)

	return nil
}

// SetCursor sets the cursor position of the view at the given point,
//
// Rules:
//   y < total lines && y > 0
//   (x < view width || x < y's line width) && x > 0
func (v *View) SetCursor(x, y int) error {
	maxX, _ := v.Size()
	if x < 0 || y < 0 || y >= len(v.lines) || (len(v.lines[y]) >= x && x >= maxX) {
		return ErrInvalidPoint
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
		return ErrInvalidPoint
	}
	v.ox = x
	v.oy = y
	return nil
}

// Origin returns the origin position of the view.
func (v *View) Origin() (x, y int) {
	return v.ox, v.oy
}

// SetWritePos sets the write position of the view's internal buffer.
// So the next Write call would write directly to the specified position.
func (v *View) SetWritePos(x, y int) error {
	if x < 0 || y < 0 {
		return ErrInvalidPoint
	}
	v.wx = x
	v.wy = y
	return nil
}

// WritePos returns the current write position of the view's internal buffer.
func (v *View) WritePos() (x, y int) {
	return v.wx, v.wy
}

// SetReadPos sets the read position of the view's internal buffer.
// So the next Read call would read from the specified position.
func (v *View) SetReadPos(x, y int) error {
	if x < 0 || y < 0 {
		return ErrInvalidPoint
	}
	v.readBuffer = nil
	v.rx = x
	v.ry = y
	return nil
}

// ReadPos returns the current read position of the view's internal buffer.
func (v *View) ReadPos() (x, y int) {
	return v.rx, v.ry
}

// makeWriteable creates empty cells if required to make position (x, y) writeable.
func (v *View) makeWriteable(x, y int) {
	// TODO: make this more efficient

	// line `y` must be index-able (that's why `<=`)
	for len(v.lines) <= y {
		if cap(v.lines) > len(v.lines) {
			newLen := cap(v.lines)
			if newLen > y {
				newLen = y + 1
			}
			v.lines = v.lines[:newLen]
		} else {
			v.lines = append(v.lines, nil)
		}
	}
	// cell `x` must not be index-able (that's why `<`)
	// append should be used by `lines[y]` user if he wants to write beyond `x`
	for len(v.lines[y]) < x {
		if cap(v.lines[y]) > len(v.lines[y]) {
			newLen := cap(v.lines[y])
			if newLen > x {
				newLen = x
			}
			v.lines[y] = v.lines[y][:newLen]
		} else {
			v.lines[y] = append(v.lines[y], cell{})
		}
	}
}

// writeCells copies []cell to specified location (x, y)
// !!! caller MUST ensure that specified location (x, y) is writeable by calling makeWriteable
func (v *View) writeCells(x, y int, cells []cell) {
	var newLen int
	// use maximum len available
	line := v.lines[y][:cap(v.lines[y])]
	maxCopy := len(line) - x
	if maxCopy < len(cells) {
		copy(line[x:], cells[:maxCopy])
		line = append(line, cells[maxCopy:]...)
		newLen = len(line)
	} else { // maxCopy >= len(cells)
		copy(line[x:], cells)
		newLen = x + len(cells)
		if newLen < len(v.lines[y]) {
			newLen = len(v.lines[y])
		}
	}
	v.lines[y] = line[:newLen]
}

// Write appends a byte slice into the view's internal buffer. Because
// View implements the io.Writer interface, it can be passed as parameter
// of functions like fmt.Fprintf, fmt.Fprintln, io.Copy, etc. Clear must
// be called to clear the view's buffer.
func (v *View) Write(p []byte) (n int, err error) {
	v.tainted = true
	v.writeMutex.Lock()
	v.makeWriteable(v.wx, v.wy)
	v.writeRunes(bytes.Runes(p))
	v.writeMutex.Unlock()

	return len(p), nil
}

func (v *View) WriteRunes(p []rune) {
	v.tainted = true

	// Fill with empty cells, if writing outside current view buffer
	v.makeWriteable(v.wx, v.wy)
	v.writeRunes(p)
}

func (v *View) WriteString(s string) {
	v.WriteRunes([]rune(s))
}

// writeRunes copies slice of runes into internal lines buffer.
// caller must make sure that writing position is accessable.
func (v *View) writeRunes(p []rune) {
	for _, r := range p {
		switch r {
		case '\n':
			v.wy++
			if v.wy >= len(v.lines) {
				v.lines = append(v.lines, nil)
			}

			fallthrough
			// not valid in every OS, but making runtime OS checks in cycle is bad.
		case '\r':
			v.wx = 0
		default:
			cells := v.parseInput(r)
			if cells == nil {
				continue
			}
			v.writeCells(v.wx, v.wy, cells)
			v.wx += len(cells)
		}
	}
}

// parseInput parses char by char the input written to the View. It returns nil
// while processing ESC sequences. Otherwise, it returns a cell slice that
// contains the processed data.
func (v *View) parseInput(ch rune) []cell {
	cells := []cell{}

	isEscape, err := v.ei.parseOne(ch)
	if err != nil {
		for _, r := range v.ei.runes() {
			c := cell{
				fgColor: v.FgColor,
				bgColor: v.BgColor,
				chr:     r,
			}
			cells = append(cells, c)
		}
		v.ei.reset()
	} else {
		if isEscape {
			return nil
		}
		repeatCount := 1
		if ch == '\t' {
			ch = ' '
			repeatCount = 4
		}
		for i := 0; i < repeatCount; i++ {
			c := cell{
				fgColor: v.ei.curFgColor,
				bgColor: v.ei.curBgColor,
				chr:     ch,
			}
			cells = append(cells, c)
		}
	}

	return cells
}

// Read reads data into p from the current reading position set by SetReadPos.
// It returns the number of bytes read into p.
// At EOF, err will be io.EOF.
func (v *View) Read(p []byte) (n int, err error) {
	buffer := make([]byte, utf8.UTFMax)
	offset := 0
	if v.readBuffer != nil {
		copy(p, v.readBuffer)
		if len(v.readBuffer) >= len(p) {
			if len(v.readBuffer) > len(p) {
				v.readBuffer = v.readBuffer[len(p):]
			}
			return len(p), nil
		}
		v.readBuffer = nil
	}
	for v.ry < len(v.lines) {
		for v.rx < len(v.lines[v.ry]) {
			count := utf8.EncodeRune(buffer, v.lines[v.ry][v.rx].chr)
			copy(p[offset:], buffer[:count])
			v.rx++
			newOffset := offset + count
			if newOffset >= len(p) {
				if newOffset > len(p) {
					v.readBuffer = buffer[newOffset-len(p):]
				}
				return len(p), nil
			}
			offset += count
		}
		v.rx = 0
		v.ry++
	}
	return offset, io.EOF
}

// Rewind sets read and write pos to (0, 0).
func (v *View) Rewind() {
	if err := v.SetReadPos(0, 0); err != nil {
		// SetReadPos returns error only if x and y are negative
		// we are passing 0, 0, thus no error should occur.
		panic(err)
	}
	if err := v.SetWritePos(0, 0); err != nil {
		// SetWritePos returns error only if x and y are negative
		// we are passing 0, 0, thus no error should occur.
		panic(err)
	}
}

// viewLines returns the lines to render on the screen
func (v *View) viewLines() [][]cell {
	if !v.Wrap {
		return v.lines
	}

	renderLines := [][]cell{}
	for _, viewLine := range v.lines {
		for {
			lineToRender, _, end := v.takeLine(&viewLine)
			renderLines = append(renderLines, lineToRender)
			if end {
				break
			}
		}
	}
	return renderLines
}

// IsTainted tells us if the view is tainted
func (v *View) IsTainted() bool {
	return v.tainted
}

// draw re-draws the view's contents.
func (v *View) draw() error {
	if !v.Visible {
		return nil
	}

	maxX, maxY := v.Size()

	if v.Wrap {
		if maxX == 0 {
			// Just return here, there is no need to try drawing chars in a too small frame
			// Nor is it needed to return an error, there is just no space
			return nil
		}
		v.ox = 0
	}

	if !v.tainted && v.contentCache != nil {
		for _, cell := range v.contentCache {
			if err := v.setRune(cell.x, cell.y, cell.chr, cell.fgColor, cell.bgColor); err != nil {
				return err
			}
		}
		return nil
	}

	linesToRender := v.viewLines()

	if v.Autoscroll && len(linesToRender) > maxY {
		v.oy = len(v.lines) - maxY
	}

	newCache := []cellCache{}
	y := 0
	for lineIndex, line := range linesToRender {
		if lineIndex < v.oy {
			continue
		}
		if y >= maxY {
			break // No need to render out of screen chars
		}

		x := 0
		for charIndex, char := range line {
			if charIndex < v.ox {
				continue
			}
			if x >= maxX {
				break // No need to render out of screen chars
			}

			fgColor := char.fgColor
			if fgColor == ColorDefault {
				fgColor = v.FgColor
			}
			bgColor := char.bgColor
			if bgColor == ColorDefault {
				bgColor = v.BgColor
			}

			newCache = append(newCache, cellCache{
				chr:     char.chr,
				bgColor: bgColor,
				fgColor: fgColor,
				x:       x,
				y:       y,
			})
			if err := v.setRune(x, y, char.chr, fgColor, bgColor); err != nil {
				return err
			}
			x += runewidth.RuneWidth(char.chr)
		}
		y++
	}

	v.contentCache = newCache
	return nil
}

// Clear empties the view's internal buffer.
// And resets reading and writing offsets.
func (v *View) Clear() {
	v.writeMutex.Lock()
	v.Rewind()
	v.tainted = true
	v.ei.reset()
	v.lines = nil
	v.clearRunes()
	v.writeMutex.Unlock()
}

// linesPosOnScreen returns based on the view lines the x and y location
// the viewX and viewY are NOT based on the view offsets
// isOnScreen is false if the selected corodinates is not on screen (this is based on the view offsets)
func (v *View) linesPosOnScreen(x, y int) (viewX int, viewY int, visable bool) {
	if x < 0 || y < 0 {
		return
	}

	maxX, maxY := v.Size()
	if !v.Wrap {
		viewX = x
		viewY = y
		visable = viewY >= v.oy && viewY < v.oy+maxY && viewX >= v.ox && viewX < v.ox+maxX
		return
	}

	var line []cell
	found := false

	for lineIndex, viewLine := range v.lines {
		if lineIndex == y {
			line = viewLine
			found = true
			break
		}

		for {
			_, _, end := v.takeLine(&viewLine)
			viewY++
			if end {
				break
			}
		}
	}

	if found {
		for {
			lineChars, width, end := v.takeLine(&line)
			lenLineChars := len(lineChars)
			if x < lenLineChars {
				x = lineWidth(lineChars[:x])
				break
			} else {
				x -= lenLineChars
			}

			if end {
				x += width
				break
			}
			viewY++
		}
	} else {
		if y < len(v.lines) {
			viewY = y
		} else {
			viewY += y - len(v.lines)
		}
	}

	viewY += x / maxX
	viewX = x - ((x / maxX) * maxX)

	visable = viewY >= v.oy && viewY < v.oy+maxY && viewX >= v.ox && viewX < v.ox+maxX
	return
}

// clearRunes erases all the cells in the view.
func (v *View) clearRunes() {
	maxX, maxY := v.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			tcellSetCell(v.x0+x+1, v.y0+y+1, ' ', v.FgColor, v.BgColor, v.outMode)
		}
	}
}

// BufferLines returns the lines in the view's internal
// buffer.
func (v *View) BufferLines() []string {
	lines := make([]string, len(v.lines))
	for i, l := range v.lines {
		str := lineType(l).String()
		str = strings.Replace(str, "\x00", " ", -1)
		lines[i] = str
	}
	return lines
}

// Buffer returns a string with the contents of the view's internal
// buffer.
func (v *View) Buffer() string {
	return linesToString(v.lines)
}

// ViewBufferLines returns the lines in the view's internal
// buffer that is shown to the user.
func (v *View) ViewBufferLines() []string {
	lines := make([]string, len(v.lines))
	for i, line := range v.lines {
		str := lineType(line).String()
		str = strings.Replace(str, "\x00", " ", -1)
		lines[i] = str
	}
	return lines
}

// LinesHeight is the count of view lines (i.e. lines excluding wrapping)
func (v *View) LinesHeight() int {
	return len(v.lines)
}

// ViewLinesHeight is the count of view lines (i.e. lines including wrapping)
func (v *View) ViewLinesHeight() int {
	if !v.tainted && v.contentCache != nil && len(v.contentCache) > 0 {
		// Use the cache if availabe, it's just a bit faster than re-calculating all frame cells
		return v.contentCache[len(v.contentCache)-1].y + 1
	}
	return len(v.viewLines())
}

// ViewBuffer returns a string with the contents of the view's buffer that is
// shown to the user.
func (v *View) ViewBuffer() string {
	return linesToString(v.lines)
}

// Line returns a string with the line of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Line(y int) (string, error) {
	if y < 0 || y >= len(v.lines) {
		return "", ErrInvalidPoint
	}

	return lineType(v.lines[y]).String(), nil
}

// Word returns a string with the word of the view's internal buffer
// at the position corresponding to the point (x, y).
func (v *View) Word(x, y int) (string, error) {
	if x < 0 || y < 0 || y >= len(v.lines) || x >= len(v.lines[y]) {
		return "", ErrInvalidPoint
	}

	str := lineType(v.lines[y]).String()

	nl := strings.LastIndexFunc(str[:x], indexFunc)
	if nl == -1 {
		nl = 0
	} else {
		nl = nl + 1
	}
	nr := strings.IndexFunc(str[x:], indexFunc)
	if nr == -1 {
		nr = len(str)
	} else {
		nr = nr + x
	}
	return string(str[nl:nr]), nil
}

// indexFunc allows to split lines by words taking into account spaces
// and 0.
func indexFunc(r rune) bool {
	return r == ' ' || r == 0
}

// SetLine changes the contents of an existing line.
func (v *View) SetLine(y int, text string) error {
	if y < 0 || y >= len(v.lines) {
		err := ErrInvalidPoint
		return err
	}

	v.tainted = true
	line := make([]cell, 0)
	for _, r := range text {
		c := v.parseInput(r)
		line = append(line, c...)
	}
	v.lines[y] = line
	return nil
}

// SetHighlight toggles highlighting of separate lines, for custom lists
// or multiple selection in views.
func (v *View) SetHighlight(y int, on bool) error {
	if y < 0 || y >= len(v.lines) {
		err := ErrInvalidPoint
		return err
	}

	line := v.lines[y]
	cells := make([]cell, 0)
	for _, c := range line {
		if on {
			c.bgColor = v.SelBgColor
			c.fgColor = v.SelFgColor
		} else {
			c.bgColor = v.BgColor
			c.fgColor = v.FgColor
		}
		cells = append(cells, c)
	}
	v.tainted = true
	v.lines[y] = cells
	return nil
}

func lineWidth(line []cell) (n int) {
	for i := range line {
		n += runewidth.RuneWidth(line[i].chr)
	}

	return
}

// takeLine slices one visable line from l and returns the sliced part
func (v *View) takeLine(l *[]cell) (visableLine []cell, width int, end bool) {
	if l == nil {
		panic("take line l can't be nil")
	}

	visableLine = []cell{}

	if len(*l) == 0 {
		end = true
		width = 0
		return
	}

	maxX, _ := v.Size()
	i := 0
	cell := cell{}

	for i, cell = range *l {
		chr := cell.chr
		charWidth := runewidth.RuneWidth(chr)

		if width+charWidth > maxX {
			break
		}

		width += charWidth
		visableLine = append(visableLine, cell)
		if width == maxX {
			break
		}
	}

	i++
	end = i == len(*l)
	*l = (*l)[i:]

	return
}

func linesToString(lines [][]cell) string {
	str := make([]string, len(lines))
	for i := range lines {
		rns := make([]rune, 0, len(lines[i]))
		line := lineType(lines[i]).String()
		for _, c := range line {
			if c != '\x00' {
				rns = append(rns, c)
			}
		}
		str[i] = string(rns)
	}

	return strings.Join(str, "\n")
}

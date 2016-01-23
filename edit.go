// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

const maxInt = int(^uint(0) >> 1)

// Edit allows to define the editor that manages the edition mode,
// including keybindings or cursor behaviour. DefaultEditor is used by
// default.
var Edit = EditorFunc(DefaultEditor)

// Editor interface must be satisfied by gocui editors.
type Editor interface {
	Edit(v *View, key Key, ch rune, mod Modifier)
}

// The EditorFunc type is an adapter to allow the use of ordinary functions as
// Editors. If f is a function with the appropriate signature, EditorFunc(f)
// is an Editor object that calls f.
type EditorFunc func(v *View, key Key, ch rune, mod Modifier)

// Edit calls f(v, key, ch, mod)
func (f EditorFunc) Edit(v *View, key Key, ch rune, mod Modifier) {
	f(v, key, ch, mod)
}

// DefaultEditor is used as the default gocui editor.
func DefaultEditor(v *View, key Key, ch rune, mod Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == KeySpace:
		v.EditWrite(' ')
	case key == KeyBackspace || key == KeyBackspace2:
		v.EditDelete(true)
	case key == KeyDelete:
		v.EditDelete(false)
	case key == KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == KeyEnter:
		v.EditNewLine()
	case key == KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

// EditWrite writes a rune at the cursor position.
func (v *View) EditWrite(ch rune) {
	v.writeRune(v.cx, v.cy, ch)
	v.MoveCursor(1, 0, true)
}

// EditDelete deletes a rune at the cursor position. back determines the
// direction.
func (v *View) EditDelete(back bool) {
	x, y := v.ox+v.cx, v.oy+v.cy
	if y < 0 {
		return
	} else if y >= len(v.viewLines) {
		v.MoveCursor(-1, 0, true)
		return
	}

	maxX, _ := v.Size()
	if back {
		if x == 0 { // start of the line
			if y < 1 {
				return
			}

			var maxPrevWidth int
			if v.Wrap {
				maxPrevWidth = maxX
			} else {
				maxPrevWidth = maxInt
			}

			if v.viewLines[y].linesX == 0 { // regular line
				v.mergeLines(v.cy - 1)
				if len(v.viewLines[y-1].line) < maxPrevWidth {
					v.MoveCursor(-1, 0, true)
				}
			} else { // wrapped line
				v.deleteRune(len(v.viewLines[y-1].line)-1, v.cy-1)
				v.MoveCursor(-1, 0, true)
			}
		} else { // middle/end of the line
			v.deleteRune(v.cx-1, v.cy)
			v.MoveCursor(-1, 0, true)
		}
	} else {
		if x == len(v.viewLines[y].line) { // end of the line
			v.mergeLines(v.cy)
		} else { // start/middle of the line
			v.deleteRune(v.cx, v.cy)
		}
	}
}

// EditNewLine inserts a new line under the cursor.
func (v *View) EditNewLine() {
	v.breakLine(v.cx, v.cy)

	y := v.oy + v.cy
	if y >= len(v.viewLines) || (y >= 0 && y < len(v.viewLines) &&
		!(v.Wrap && v.cx == 0 && v.viewLines[y].linesX > 0)) {
		// new line at the end of the buffer or
		// cursor is not at the beginning of a wrapped line
		v.ox = 0
		v.cx = 0
		v.MoveCursor(0, 1, true)
	}
}

// MoveCursor moves the cursor taking into account the width of the line/view,
// displacing the origin if necessary.
func (v *View) MoveCursor(dx, dy int, writeMode bool) {
	maxX, maxY := v.Size()
	cx, cy := v.cx+dx, v.cy+dy
	x, y := v.ox+cx, v.oy+cy

	var curLineWidth, prevLineWidth int
	// get the width of the current line
	if writeMode {
		if v.Wrap {
			curLineWidth = maxX - 1
		} else {
			curLineWidth = maxInt
		}
	} else {
		if y >= 0 && y < len(v.viewLines) {
			curLineWidth = len(v.viewLines[y].line)
			if v.Wrap && curLineWidth >= maxX {
				curLineWidth = maxX - 1
			}
		} else {
			curLineWidth = 0
		}
	}
	// get the width of the previous line
	if y-1 >= 0 && y-1 < len(v.viewLines) {
		prevLineWidth = len(v.viewLines[y-1].line)
	} else {
		prevLineWidth = 0
	}

	// adjust cursor's x position and view's x origin
	if x > curLineWidth { // move to next line
		if dx > 0 { // horizontal movement
			if !v.Wrap {
				v.ox = 0
			}
			v.cx = 0
			cy++
		} else { // vertical movement
			if curLineWidth > 0 { // move cursor to the EOL
				if v.Wrap {
					v.cx = curLineWidth
				} else {
					ncx := curLineWidth - v.ox
					if ncx < 0 {
						v.ox += ncx
						if v.ox < 0 {
							v.ox = 0
						}
						v.cx = 0
					} else {
						v.cx = ncx
					}
				}
			} else {
				if !v.Wrap {
					v.ox = 0
				}
				v.cx = 0
			}
		}
	} else if cx < 0 {
		if !v.Wrap && v.ox > 0 { // move origin to the left
			v.ox--
		} else { // move to previous line
			if prevLineWidth > 0 {
				if !v.Wrap { // set origin so the EOL is visible
					nox := prevLineWidth - maxX + 1
					if nox < 0 {
						v.ox = 0
					} else {
						v.ox = nox
					}
				}
				v.cx = prevLineWidth
			} else {
				if !v.Wrap {
					v.ox = 0
				}
				v.cx = 0
			}
			cy--
		}
	} else { // stay on the same line
		if v.Wrap {
			v.cx = cx
		} else {
			if cx >= maxX {
				v.ox++
			} else {
				v.cx = cx
			}
		}
	}

	// adjust cursor's y position and view's y origin
	if cy >= maxY {
		v.oy++
	} else if cy < 0 {
		if v.oy > 0 {
			v.oy--
		}
	} else {
		v.cy = cy
	}
}

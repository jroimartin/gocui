// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/nsf/termbox-go"

const MaxInt = int(^uint(0) >> 1)

// handleEdit manages the edition mode. We do not handle errors here because if
// an error happens, it is enough to keep the view without modifications.
func (g *Gui) handleEdit(v *View, ev *termbox.Event) {
	switch {
	case ev.Ch != 0 && ev.Mod == 0:
		v.editWrite(ev.Ch)
	case ev.Key == termbox.KeySpace:
		v.editWrite(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		v.editDelete(true)
	case ev.Key == termbox.KeyDelete:
		v.editDelete(false)
	case ev.Key == termbox.KeyInsert:
		v.overwrite = !v.overwrite
	case ev.Key == termbox.KeyEnter:
		v.editNewLine()
	case ev.Key == termbox.KeyArrowDown:
		v.moveCursor(0, 1, false)
	case ev.Key == termbox.KeyArrowUp:
		v.moveCursor(0, -1, false)
	case ev.Key == termbox.KeyArrowLeft:
		v.moveCursor(-1, 0, false)
	case ev.Key == termbox.KeyArrowRight:
		v.moveCursor(1, 0, false)
	}
}

// editWrite writes a rune at the cursor position.
func (v *View) editWrite(ch rune) {
	v.writeRune(v.cx, v.cy, ch)
	v.moveCursor(1, 0, true)
}

// editDelete deletes a rune at the cursor position. back determines
// the direction.
func (v *View) editDelete(back bool) {
	if back {
		if v.cx == 0 {
			v.mergeLines(v.cy - 1)
		} else {
			v.deleteRune(v.cx-1, v.cy)
		}
		v.moveCursor(-1, 0, true)
	} else {
		y := v.oy + v.cy
		if y >= 0 && y < len(v.viewLines) && v.cx == len(v.viewLines[y].line) {
			v.mergeLines(v.cy)
		} else {
			v.deleteRune(v.cx, v.cy)
		}
	}
}

// editNewLine inserts a new line under the cursor.
func (v *View) editNewLine() {
	v.breakLine(v.cx, v.cy)

	y := v.oy + v.cy
	if y >= len(v.viewLines) || (y >= 0 && y < len(v.viewLines) &&
		!(v.Wrap && v.cx == 0 && v.viewLines[y].linesX > 0)) {
		// new line at the end of the buffer or
		// cursor is not at the beginning of a wrapped line
		v.ox = 0
		v.cx = 0
		v.moveCursor(0, 1, true)
	}
}

// moveCursor moves the cursor taking into account the line or view widths and
// moves the origin when necessary. If writeMode is false, the cursor jumps to
// the next line when it reaches the end of the line, otherwise it jumps when
// the cursor reaches the width of the view.
func (v *View) moveCursor(dx, dy int, writeMode bool) {
	maxX, maxY := v.Size()
	cx, cy := v.cx+dx, v.cy+dy
	x, y := v.ox+cx, v.oy+cy

	var curLineWidth, prevLineWidth int
	// get the width of the current line
	if writeMode {
		if v.Wrap {
			curLineWidth = maxX - 1
		} else {
			curLineWidth = MaxInt
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
			cy += 1
		} else { // vertical movement
			if curLineWidth > 0 { // move cursor to the EOF
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
			v.ox -= 1
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
			cy -= 1
		}
	} else { // stay on the same line
		if v.Wrap {
			v.cx = cx
		} else {
			if cx >= maxX {
				v.ox += 1
			} else {
				v.cx = cx
			}
		}
	}

	// adjust cursor's y position and view's y origin
	if cy >= maxY {
		v.oy += 1
	} else if cy < 0 {
		if v.oy > 0 {
			v.oy -= 1
		}
	} else {
		v.cy = cy
	}
}

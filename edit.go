// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

// editWrite writes a rune in edit mode.
func (v *View) editWrite(ch rune) error {
	maxX, _ := v.Size()
	v.writeRune(v.cx, v.cy, ch)
	if v.cx == maxX-1 {
		if err := v.SetOrigin(v.ox+1, v.oy); err != nil {
			return err
		}
	} else {
		if err := v.SetCursor(v.cx+1, v.cy); err != nil {
			return err
		}
	}
	return nil
}

// editDelete deletes a rune in edit mode. back determines the direction.
func (v *View) editDelete(back bool) error {
	if back {
		v.deleteRune(v.cx-1, v.cy)
		if v.cx == 0 {
			if v.ox > 0 {
				if err := v.SetOrigin(v.ox-1, v.oy); err != nil {
					return err
				}
			}
		} else {
			if err := v.SetCursor(v.cx-1, v.cy); err != nil {
				return err
			}
		}
	} else {
		v.deleteRune(v.cx, v.cy)
	}
	return nil
}

// editLine inserts a new line under the cursor in edit mode.
func (v *View) editLine() error {
	_, maxY := v.Size()
	v.addLine(v.cy + 1)
	if v.cy == maxY-1 {
		if err := v.SetOrigin(0, v.oy+1); err != nil {
			return err
		}
		if err := v.SetCursor(0, v.cy); err != nil {
			return err
		}
	} else {
		if err := v.SetCursor(0, v.cy+1); err != nil {
			return err
		}
	}
	return nil
}

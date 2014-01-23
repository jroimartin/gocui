// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

// editWrite writes a rune at the cursor position.
func (v *View) editWrite(ch rune) error {
	v.writeRune(v.cx, v.cy, ch)
	if err := v.SetCursor(v.cx+1, v.cy); err != nil {
		if err := v.SetOrigin(v.ox+1, v.oy); err != nil {
			return err
		}
	}
	return nil
}

// editDelete deletes a rune at the cursor position. back determines
// the direction.
func (v *View) editDelete(back bool) error {
	if back {
		v.deleteRune(v.cx-1, v.cy)
		if err := v.SetCursor(v.cx-1, v.cy); err != nil && v.ox > 0 {
			if err := v.SetOrigin(v.ox-1, v.oy); err != nil {
				return err
			}
		}
	} else {
		v.deleteRune(v.cx, v.cy)
	}
	return nil
}

// editLine inserts a new line under the cursor.
func (v *View) editLine() error {
	v.addLine(v.cy + 1)
	if err := v.SetCursor(v.cx, v.cy+1); err != nil {
		if err := v.SetOrigin(v.ox, v.oy+1); err != nil {
			return err
		}
	}
	if err := v.SetCursor(0, v.cy); err != nil {
		return err
	}
	if err := v.SetOrigin(0, v.oy); err != nil {
		return err
	}
	return nil
}

// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

// editWrite writes a rune at the cursor position.
func (v *View) editWrite(ch rune) error {
	v.writeRune(v.cx, v.cy, ch)
	_ = v.getRuneLen(ch)
	//if err := v.SetCursor(v.cx+l, v.cy); err != nil {
	/*if err := v.SetOrigin(v.ox+l, v.oy); err != nil {
		return err
	}*/
	//}
	return nil
}

// editDelete deletes a rune at the cursor position. back determines
// the direction.
func (v *View) editDelete(back bool) error {
	if back {
		_ = v.deleteRune(v.cx-1, v.cy)
		//if err := v.SetCursor(v.cx-l, v.cy); err != nil && v.ox > 0 {
		/*if err := v.SetOrigin(v.ox-l, v.oy); err != nil {
			return err
		}*/
		//}
	} else {
		v.deleteRune(v.cx, v.cy)
	}
	return nil
}

// editLine inserts a new line under the cursor.
func (v *View) editLine() error {
	v.AddLine(v.cx, v.cy)
	return nil
}

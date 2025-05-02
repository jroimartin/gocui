//Help editor to Input cjk language.
//Author fuhuizn@163.com
package gocui

import (
	"bytes"

	"github.com/mattn/go-runewidth"
)

func modifyCJK(p []byte) []byte {
	buf := bytes.NewBuffer(bytes.Trim(p, " \n\t"))
	sz := len(buf.String())
	buf1 := bytes.NewBufferString("")
	var r rune
	var wr bool
	for i := 0; i < sz; i++ {
		r, _, _ = buf.ReadRune()
		if r != rune(0) && wr == false {
			buf1.WriteRune(r)
		} else if wr == true {
			if r != rune(' ') {
				buf1.WriteRune(r)
			}
		}
		wr = runewidth.RuneWidth(r) > 1
	}
	return buf1.Bytes()
}

//ReadEditor Read byte array from editor 'v', delete the auto appended blank after CJK runes.
func (v *View) ReadEditor() []byte {
	var b = make([]byte, 300)
	n, _ := v.Read(b)
	if n > 0 {
		return modifyCJK(b[:n])
	} else {
		return nil
	}
}

// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/nsf/termbox-go"
)

type View struct {
	Name                   string
	X0, Y0, X1, Y1         int
	CX, CY                 int
	OX, OY                 int
	Highlight              bool
	buffer                 []rune
	bgColor, fgColor       Attribute
	selBgColor, selFgColor Attribute
}

func newView(name string, x0, y0, x1, y1 int) *View {
	v := &View{
		Name: name,
		X0:   x0,
		Y0:   y0,
		X1:   x1,
		Y1:   y1,
	}
	return v
}

func (v *View) Size() (x, y int) {
	return v.X1 - v.X0 - 1, v.Y1 - v.Y0 - 1
}

func (v *View) setRune(x, y int, ch rune) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}

	var fgColor, bgColor Attribute
	if v.Highlight && y == v.CY {
		fgColor = v.selFgColor
		bgColor = v.selBgColor
	} else {
		fgColor = v.fgColor
		bgColor = v.bgColor
	}
	termbox.SetCell(v.X0+x+1, v.Y0+y+1, ch,
		termbox.Attribute(fgColor), termbox.Attribute(bgColor))
	return nil
}

func (v *View) SetCursor(x, y int) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}
	v.CX = x
	v.CY = y
	return nil
}

func (v *View) SetOrigin(x, y int) {
	v.OX = x
	v.OY = y
}

func (v *View) Write(p []byte) (n int, err error) {
	pr := bytes.Runes(p)
	v.buffer = append(v.buffer, pr...)
	return len(pr), nil
}

func (v *View) draw() error {
	maxX, maxY := v.Size()
	buf := bytes.NewBufferString(string(v.buffer))
	br := bufio.NewReader(buf)

	y := 0
	for i := 0; ; i++ {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if i < v.OY {
			continue
		}
		x := 0
		for j, ch := range bytes.Runes(line) {
			if j < v.OX {
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

func (v *View) Clear() {
	v.buffer = nil
	v.clearRunes()
}

func (v *View) clearRunes() {
	maxX, maxY := v.Size()
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			termbox.SetCell(v.X0+x+1, v.Y0+y+1, 0,
				termbox.Attribute(v.fgColor), termbox.Attribute(v.bgColor))
		}
	}
}

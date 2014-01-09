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
	buffer                 []rune
	X0, Y0, X1, Y1         int
	CX, CY                 int
	BgColor, FgColor       termbox.Attribute
	SelBgColor, SelFgColor termbox.Attribute
}

func NewView(name string, x0, y0, x1, y1 int) (v *View) {
	v = &View{
		Name:       name,
		X0:         x0,
		Y0:         y0,
		X1:         x1,
		Y1:         y1,
		BgColor:    termbox.ColorBlack,
		FgColor:    termbox.ColorWhite,
		SelBgColor: termbox.ColorBlack,
		SelFgColor: termbox.ColorWhite,
	}
	return v
}

func (v *View) Size() (x, y int) {
	return v.X1 - v.X0 - 1, v.Y1 - v.Y0 - 1
}

func (v *View) SetRune(x, y int, ch rune) (err error) {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}
	termbox.SetCell(v.X0+x+1, v.Y0+y+1, ch, v.FgColor, v.BgColor)
	return nil
}

func (v *View) GetRune(x, y int) (ch rune, err error) {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return 0, errors.New("invalid point")
	}
	h := v.Y1 - v.Y0 - 1
	c := v.buffer[y*h+x]
	return c, nil
}

func (v *View) SetCursor(x, y int) (err error) {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return errors.New("invalid point")
	}
	v.CX = x
	v.CY = y
	return nil
}

func (v *View) Write(p []byte) (n int, err error) {
	pr := bytes.Runes(p)
	v.buffer = append(v.buffer, pr...)
	return len(pr), nil
}

func (v *View) Draw() (err error) {
	maxX, maxY := v.Size()
	buf := bytes.NewBufferString(string(v.buffer))
	br := bufio.NewReader(buf)
	for nl := 0; ; nl++ {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		for i, ch := range bytes.Runes(line) {
			if i >= 0 && i < maxX && nl >= 0 && nl < maxY {
				if err := v.SetRune(i, nl, ch); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

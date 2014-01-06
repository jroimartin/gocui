package gocui

import (
	"errors"
	"github.com/nsf/termbox-go"
)

type View struct {
	Name                   string
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

func (v *View) SetCursor(x, y int) (err error) {
	if x < 0 || v.X0+x+1 >= v.X1 || y < 0 || v.Y0+y+1 >= v.Y1 {
		return errors.New("invalid point")
	}
	v.CX = x
	v.CY = y
	return nil
}

func (v *View) Write(p []byte) (n int, err error) {
	return 0, nil
}

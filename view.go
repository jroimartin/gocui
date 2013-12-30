package gocui

import (
	"github.com/jroimartin/termbox-go"
)

type View struct {
	Name                   string
	X0, Y0, X1, Y1         int
	cx, cy                 int
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

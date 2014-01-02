package gocui

import (
	"github.com/nsf/termbox-go"
)

type View struct {
	Name                   string
	X, Y, W, H             float32
	X0, Y0, X1, Y1         int
	cx, cy                 int
	BgColor, FgColor       termbox.Attribute
	SelBgColor, SelFgColor termbox.Attribute
}

func NewView(name string, x, y, w, h float32, maxX, maxY int) (v *View) {
	v = &View{
		Name:       name,
		X:          x,
		Y:          y,
		W:          w,
		H:          h,
		BgColor:    termbox.ColorBlack,
		FgColor:    termbox.ColorWhite,
		SelBgColor: termbox.ColorBlack,
		SelFgColor: termbox.ColorWhite,
	}
	v.Resize(maxX, maxY)
	return v
}

func (v *View) Resize(maxX, maxY int) {
	switch {
	case v.X0 < 0:
		v.X0 = -1
	case v.X0 < 1:
		v.X0 = int(v.X*float32(maxX) + 0.5)
	default:
		v.X0 = int(v.X + 0.5)
	}
	switch {
	case v.W < 0:
		v.X1 = maxX
	case v.W < 1:
		v.X1 = v.X0 + int(v.W*float32(maxX)+0.5)
	default:
		v.X1 = v.X0 + int(v.W+0.5)
	}

	switch {
	case v.Y < 0:
		v.Y0 = -1
	case v.Y < 1:
		v.Y0 = int(v.Y*float32(maxY) + 0.5)
	default:
		v.Y0 = int(v.Y + 0.5)
	}
	switch {
	case v.H < 0:
		v.X1 = maxY
	case v.H < 1:
		v.Y1 = v.Y0 + int(v.H*float32(maxY)+0.5)
	default:
		v.Y1 = v.Y0 + int(v.H+0.5)
	}
}

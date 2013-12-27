package gocui

import (
	"github.com/jroimartin/termbox-go"
)

type View struct {
	x0, y0, x1, y1         int
	cx, cy                 int
	BgColor, FgColor       termbox.Attribute
	SelBgColor, SelFgColor termbox.Attribute
}

func NewView(x0, y0, x1, y1 int) (v *View) {
	return &View{x0: x0, y0: y0, x1: x1, y1: y1,
		BgColor: termbox.ColorBlack, FgColor: termbox.ColorWhite,
		SelBgColor: termbox.ColorBlack, SelFgColor: termbox.ColorWhite}
}

func (v *View) SetCell(x, y int, ch rune) {
	var bgColor, fgColor termbox.Attribute

	if y == v.cy {
		bgColor = v.SelBgColor
		fgColor = v.SelFgColor
	} else {
		bgColor = v.BgColor
		fgColor = v.FgColor
	}
	termbox.SetCell(x, y, ch, fgColor, bgColor)
}

func (v *View) Draw() (err error) {
	//maxX, maxY := termbox.Size()
	v.SetCell(v.x0, v.y0, RuneCornerTopLeft)
	v.SetCell(v.x1, v.y0, RuneCornerTopRight)
	v.SetCell(v.x0, v.y1, RuneCornerBottomLeft)
	v.SetCell(v.x1, v.y1, RuneCornerBottomRight)
	for x := v.x0 + 1; x < v.x1; x++ {
		v.SetCell(x, v.y0, RuneEdgeHorizontal)
		v.SetCell(x, v.y1, RuneEdgeHorizontal)
	}
	for y := v.y0 + 1; y < v.y1; y++ {
		v.SetCell(v.x0, y, RuneEdgeVertical)
		v.SetCell(v.x1, y, RuneEdgeVertical)
	}
	return nil
}

func (v *View) Resize() (err error) {
	//TODO
	return nil
}

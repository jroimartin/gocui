package gocui

import (
	"github.com/jroimartin/termbox-go"
)

type View struct {
	name                   string
	x0, y0, x1, y1         int
	cx, cy                 int
	BgColor, FgColor       termbox.Attribute
	SelBgColor, SelFgColor termbox.Attribute
}

func NewView(name string, x0, y0, x1, y1 int) (v *View) {
	v = &View{
		name:       name,
		x0:         x0,
		y0:         y0,
		x1:         x1,
		y1:         y1,
		BgColor:    termbox.ColorBlack,
		FgColor:    termbox.ColorWhite,
		SelBgColor: termbox.ColorBlack,
		SelFgColor: termbox.ColorWhite,
	}
	return v
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
	maxX, maxY := termbox.Size()
	if v.y0 != -1 {
		if v.x0 != -1 {
			v.SetCell(v.x0, v.y0, RuneCornerTopLeft)
		}
		if v.x1 != maxX {
			v.SetCell(v.x1, v.y0, RuneCornerTopRight)
		}
	}
	if v.y0 != maxY {
		if v.x0 != -1 {
			v.SetCell(v.x0, v.y1, RuneCornerBottomLeft)
		}
		if v.x1 != maxX {
			v.SetCell(v.x1, v.y1, RuneCornerBottomRight)
		}
	}
	for x := v.x0 + 1; x < v.x1; x++ {
		if v.y0 != -1 {
			v.SetCell(x, v.y0, RuneEdgeHorizontal)
		}
		if v.y1 != maxY {
			v.SetCell(x, v.y1, RuneEdgeHorizontal)
		}
	}
	for y := v.y0 + 1; y < v.y1; y++ {
		if v.x0 != -1 {
			v.SetCell(v.x0, y, RuneEdgeVertical)
		}
		if v.x1 != maxX {
			v.SetCell(v.x1, y, RuneEdgeVertical)
		}
	}
	return nil
}

func (v *View) Resize() (err error) {
	//TODO
	return nil
}

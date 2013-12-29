package gocui

import (
	"errors"

	"github.com/jroimartin/termbox-go"
)

var ErrorQuit = errors.New("quit")

const (
	RuneCornerTopLeft      = '┌'
	RuneCornerTopRight     = '┐'
	RuneCornerBottomLeft   = '└'
	RuneCornerBottomRight  = '┘'
	RuneSideVertical       = '│'
	RuneSideHorizontal     = '─'
	RuneIntersection       = '┼'
	RuneIntersectionLeft   = '├'
	RuneIntersectionRight  = '┤'
	RuneIntersectionTop    = '┬'
	RuneIntersectionBottom = '┴'
	RuneTriangleLeft       = '◄'
	RuneTriangleRight      = '►'
	RuneArrowLeft          = '←'
	RuneArrowRight         = '→'
)

type Gui struct {
	events           chan termbox.Event
	views            []*View
	currentView      *View
	BgColor, FgColor termbox.Attribute
}

func NewGui() (g *Gui) {
	return &Gui{}
}

func (g *Gui) Init() (err error) {
	g.events = make(chan termbox.Event, 20)
	g.BgColor = termbox.ColorWhite
	g.FgColor = termbox.ColorBlack
	return termbox.Init()
}

func (g *Gui) Close() {
	termbox.Close()
}

func (g *Gui) Size() (x, y int) {
	return termbox.Size()
}

func (g *Gui) AddView(name string, x0, y0, x1, y1 int) (v *View, err error) {
	maxX, maxY := termbox.Size()

	if x0 < -1 || y0 < -1 || x1 < -1 || y1 < -1 ||
		x0 > maxX || y0 > maxY || x1 > maxX || y1 > maxY ||
		x0 >= x1 || y0 >= y1 {
		return nil, errors.New("invalid points")
	}

	for _, v := range g.views {
		if name == v.Name {
			return nil, errors.New("invalid name")
		}
	}

	v = NewView(name, x0, y0, x1, y1)
	g.views = append(g.views, v)
	return v, nil
}

func (g *Gui) MainLoop() (err error) {
	go func() {
		for {
			g.events <- termbox.PollEvent()
		}
	}()

	if err := g.resize(); err != nil {
		return err
	}
	if err := g.draw(); err != nil {
		return err
	}
	// TODO: Set initial cursor position
	//termbox.SetCursor(10, 10)
	termbox.Flush()

	for {
		select {
		case ev := <-g.events:
			if err := g.handleEvent(&ev); err != nil {
				return err
			}
			if err := g.consumeevents(); err != nil {
				return err
			}
			if err := g.draw(); err != nil {
				return err
			}
			termbox.Flush()
		}
	}
	return nil
}

func (g *Gui) SetCell(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, g.FgColor, g.BgColor)
}

func (g *Gui) GetCell(x, y int) (ch rune, err error) {
	maxX, maxY := termbox.Size()
	if x < 0 || y < 0 || x >= maxX || y >= maxY {
		return 0, errors.New("invalid point")
	}
	c := termbox.CellBuffer()[y*maxX+x]
	return c.Ch, nil
}

func (g *Gui) consumeevents() (err error) {
	for {
		select {
		case ev := <-g.events:
			if err := g.handleEvent(&ev); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (g *Gui) handleEvent(ev *termbox.Event) (err error) {
	switch ev.Type {
	case termbox.EventKey:
		return g.onKey(ev)
	case termbox.EventResize:
		return g.resize()
	case termbox.EventError:
		return ev.Err
	default:
		return nil
	}
}

func (g *Gui) draw() (err error) {
	for _, v := range g.views {
		if err := g.drawView(v); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gui) drawView(v *View) (err error) {
	return nil
}

func (g *Gui) resize() (err error) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err := g.resizeView(); err != nil {
		return err
	}
	if err := g.drawFrames(); err != nil {
		return err
	}
	if err := g.drawFrameIntersections(); err != nil {
		return err
	}
	return nil

}

func (g *Gui) drawFrames() (err error) {
	maxX, maxY := termbox.Size()
	for _, v := range g.views {
		if v.Y0 != -1 {
			if v.X0 != -1 {
				g.SetCell(v.X0, v.Y0, RuneCornerTopLeft)
			}
			if v.X1 != maxX {
				g.SetCell(v.X1, v.Y0, RuneCornerTopRight)
			}
		}
		if v.Y1 != maxY {
			if v.X0 != -1 {
				g.SetCell(v.X0, v.Y1, RuneCornerBottomLeft)
			}
			if v.X1 != maxX {
				g.SetCell(v.X1, v.Y1, RuneCornerBottomRight)
			}
		}
		for x := v.X0 + 1; x < v.X1; x++ {
			if v.Y0 != -1 {
				g.SetCell(x, v.Y0, RuneSideHorizontal)
			}
			if v.Y1 != maxY {
				g.SetCell(x, v.Y1, RuneSideHorizontal)
			}
		}
		for y := v.Y0 + 1; y < v.Y1; y++ {
			if v.X0 != -1 {
				g.SetCell(v.X0, y, RuneSideVertical)
			}
			if v.X1 != maxX {
				g.SetCell(v.X1, y, RuneSideVertical)
			}
		}
	}
	return nil
}

func (g *Gui) drawFrameIntersections() (err error) {
	for _, v := range g.views {
		// ┌
		if ch, err := g.GetCell(v.X0, v.Y0); err == nil {
			switch ch {
			case RuneCornerTopLeft: // '┌'
				// Nothing
			case RuneCornerTopRight: // '┐'
				g.SetCell(v.X0, v.Y0, RuneIntersectionTop)
			case RuneCornerBottomLeft: // '└'
				g.SetCell(v.X0, v.Y0, RuneIntersectionLeft)
			case RuneCornerBottomRight: // '┘'
				g.SetCell(v.X0, v.Y0, RuneIntersection)
			case RuneSideVertical: // '│'
				g.SetCell(v.X0, v.Y0, RuneIntersectionLeft)
			case RuneSideHorizontal: // '─'
				g.SetCell(v.X0, v.Y0, RuneIntersectionTop)
			case RuneIntersection: // '┼'
				// Nothing
			case RuneIntersectionLeft: // '├'
				// Nothing
			case RuneIntersectionRight: // '┤'
				g.SetCell(v.X0, v.Y0, RuneIntersection)
			case RuneIntersectionTop: // '┬'
				// Nothing
			case RuneIntersectionBottom: // '┴'
				g.SetCell(v.X0, v.Y0, RuneIntersection)
			}
		}

		// ┐
		if ch, err := g.GetCell(v.X1, v.Y0); err == nil {
			switch ch {
			case RuneCornerTopLeft: // '┌'
				g.SetCell(v.X1, v.Y0, RuneIntersectionTop)
			case RuneCornerTopRight: // '┐'
				// Nothing
			case RuneCornerBottomLeft: // '└'
				g.SetCell(v.X1, v.Y0, RuneIntersection)
			case RuneCornerBottomRight: // '┘'
				g.SetCell(v.X1, v.Y0, RuneIntersectionRight)
			case RuneSideVertical: // '│'
				g.SetCell(v.X1, v.Y0, RuneIntersectionRight)
			case RuneSideHorizontal: // '─'
				g.SetCell(v.X1, v.Y0, RuneIntersectionTop)
			case RuneIntersection: // '┼'
				// Nothing
			case RuneIntersectionLeft: // '├'
				g.SetCell(v.X1, v.Y0, RuneIntersection)
			case RuneIntersectionRight: // '┤'
				// Nothing
			case RuneIntersectionTop: // '┬'
				// Nothing
			case RuneIntersectionBottom: // '┴'
				g.SetCell(v.X1, v.Y0, RuneIntersection)
			}
		}

		// └
		if ch, err := g.GetCell(v.X0, v.Y1); err == nil {
			switch ch {
			case RuneCornerTopLeft: // '┌'
				g.SetCell(v.X0, v.Y1, RuneIntersectionLeft)
			case RuneCornerTopRight: // '┐'
				g.SetCell(v.X0, v.Y1, RuneIntersection)
			case RuneCornerBottomLeft: // '└'
				// Nothing
			case RuneCornerBottomRight: // '┘'
				g.SetCell(v.X0, v.Y1, RuneIntersectionBottom)
			case RuneSideVertical: // '│'
				g.SetCell(v.X0, v.Y1, RuneIntersectionLeft)
			case RuneSideHorizontal: // '─'
				g.SetCell(v.X0, v.Y1, RuneIntersectionBottom)
			case RuneIntersection: // '┼'
				// Nothing
			case RuneIntersectionLeft: // '├'
				// Nothing
			case RuneIntersectionRight: // '┤'
				g.SetCell(v.X0, v.Y1, RuneIntersection)
			case RuneIntersectionTop: // '┬'
				g.SetCell(v.X0, v.Y1, RuneIntersection)
			case RuneIntersectionBottom: // '┴'
				// Nothing
			}
		}

		// ┘
		if ch, err := g.GetCell(v.X1, v.Y1); err == nil {
			switch ch {
			case RuneCornerTopLeft: // '┌'
				g.SetCell(v.X1, v.Y1, RuneIntersection)
			case RuneCornerTopRight: // '┐'
				g.SetCell(v.X1, v.Y1, RuneIntersectionRight)
			case RuneCornerBottomLeft: // '└'
				g.SetCell(v.X1, v.Y1, RuneIntersectionBottom)
			case RuneCornerBottomRight: // '┘'
				// Nothing
			case RuneSideVertical: // '│'
				g.SetCell(v.X1, v.Y1, RuneIntersectionRight)
			case RuneSideHorizontal: // '─'
				g.SetCell(v.X1, v.Y1, RuneIntersectionBottom)
			case RuneIntersection: // '┼'
				// Nothing
			case RuneIntersectionLeft: // '├'
				g.SetCell(v.X1, v.Y1, RuneIntersection)
			case RuneIntersectionRight: // '┤'
				// Nothing
			case RuneIntersectionTop: // '┬'
				g.SetCell(v.X1, v.Y1, RuneIntersection)
			case RuneIntersectionBottom: // '┴'
				// Nothing
			}
		}
	}
	return nil
}

func (g *Gui) resizeView() (err error) {
	return nil
}

func (g *Gui) onKey(ev *termbox.Event) (err error) {
	switch ev.Key {
	case termbox.KeyCtrlC:
		return ErrorQuit
	default:
		return nil
	}
}

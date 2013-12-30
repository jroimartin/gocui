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

func (g *Gui) SetCell(x, y int, ch rune) (err error) {
	maxX, maxY := termbox.Size()
	if x < 0 || y < 0 || x >= maxX || y >= maxY {
		return errors.New("invalid point")
	}
	termbox.SetCell(x, y, ch, g.FgColor, g.BgColor)
	return nil
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
	if err := g.resizeViews(); err != nil {
		return err
	}
	if err := g.drawFrames(); err != nil {
		return err
	}
	if err := g.drawIntersections(); err != nil {
		return err
	}
	return nil

}

func (g *Gui) drawFrames() (err error) {
	maxX, maxY := termbox.Size()
	for _, v := range g.views {
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

func (g *Gui) drawIntersections() (err error) {
	for _, v := range g.views {
		if ch, ok := g.getIntersectionRune(v.X0, v.Y0); ok {
			g.SetCell(v.X0, v.Y0, ch)
		}
		if ch, ok := g.getIntersectionRune(v.X0, v.Y1); ok {
			g.SetCell(v.X0, v.Y1, ch)
		}
		if ch, ok := g.getIntersectionRune(v.X1, v.Y0); ok {
			g.SetCell(v.X1, v.Y0, ch)
		}
		if ch, ok := g.getIntersectionRune(v.X1, v.Y1); ok {
			g.SetCell(v.X1, v.Y1, ch)
		}
	}
	return nil
}

func (g *Gui) getIntersectionRune(x, y int) (ch rune, ok bool) {
	maxX, maxY := termbox.Size()
	if x < 0 || y < 0 || x >= maxX || y >= maxY {
		return 0, false
	}

	_chTop, _ := g.GetCell(x, y-1)
	chTop := verticalRune(_chTop)
	_chBottom, _ := g.GetCell(x, y+1)
	chBottom := verticalRune(_chBottom)
	_chLeft, _ := g.GetCell(x-1, y)
	chLeft := horizontalRune(_chLeft)
	_chRight, _ := g.GetCell(x+1, y)
	chRight := horizontalRune(_chRight)

	switch {
	case !chTop && chBottom && !chLeft && chRight:
		ch = RuneCornerTopLeft // '┌'
	case !chTop && chBottom && chLeft && !chRight:
		ch = RuneCornerTopRight // '┐'
	case chTop && !chBottom && !chLeft && chRight:
		ch = RuneCornerBottomLeft // '└'
	case chTop && !chBottom && chLeft && !chRight:
		ch = RuneCornerBottomRight // '┘'
	case chTop && chBottom && chLeft && chRight:
		ch = RuneIntersection // '┼'
	case chTop && chBottom && !chLeft && chRight:
		ch = RuneIntersectionLeft // '├'
	case chTop && chBottom && chLeft && !chRight:
		ch = RuneIntersectionRight // '┤'
	case !chTop && chBottom && chLeft && chRight:
		ch = RuneIntersectionTop // '┬'
	case chTop && !chBottom && chLeft && chRight:
		ch = RuneIntersectionBottom // '┴'
	default:
		return 0, false
	}
	return ch, true
}

func verticalRune(ch rune) bool {
	if ch == RuneSideVertical || ch == RuneIntersection ||
		ch == RuneIntersectionLeft || ch == RuneIntersectionRight {
		return true
	}
	return false
}

func horizontalRune(ch rune) bool {
	if ch == RuneSideHorizontal || ch == RuneIntersection ||
		ch == RuneIntersectionTop || ch == RuneIntersectionBottom {
		return true
	}
	return false
}

func (g *Gui) resizeViews() (err error) {
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

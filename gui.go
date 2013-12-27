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
	RuneEdgeVertical       = '│'
	RuneEdgeHorizontal     = '─'
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
	events      chan termbox.Event
	views       []*View
	currentView *View
}

func NewGui() (g *Gui) {
	return &Gui{}
}

func (g *Gui) Init() (err error) {
	g.events = make(chan termbox.Event, 20)
	return termbox.Init()
}

func (g *Gui) Close() {
	termbox.Close()
}

func (g *Gui) Size() (x, y int) {
	return termbox.Size()
}

func (g *Gui) AddView(x0, y0, x1, y1 int) (v *View, err error) {
	maxX, maxY := termbox.Size()

	if x0 < 0 || y0 < 0 || x1 < 0 || y1 < 0 ||
		x0 > maxX || y0 > maxY || x1 > maxX || y1 > maxY {
		return nil, errors.New("Invalid coordinates")
	}
	v = NewView(x0, y0, x1, y1)
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
		if err := v.Draw(); err != nil {
			return err
		}
	}
	return nil

}

func (g *Gui) resize() (err error) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, v := range g.views {
		if err := v.Resize(); err != nil {
			return err
		}
	}
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

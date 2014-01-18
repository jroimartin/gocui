// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var (
	ErrorQuit    error = errors.New("quit")
	ErrorUnkView error = errors.New("unknown view")
)

type Gui struct {
	events      chan termbox.Event
	views       []*View
	currentView *View
	layout      func(*Gui) error
	keybindings []*keybinding
	maxX, maxY  int

	BgColor, FgColor       Attribute
	SelBgColor, SelFgColor Attribute
	ShowCursor             bool
}

func NewGui() *Gui {
	return &Gui{}
}

func (g *Gui) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	g.events = make(chan termbox.Event, 20)
	g.maxX, g.maxY = termbox.Size()
	g.BgColor = ColorBlack
	g.FgColor = ColorWhite
	return nil
}

func (g *Gui) Close() {
	termbox.Close()
}

func (g *Gui) Size() (x, y int) {
	return g.maxX, g.maxY
}

func (g *Gui) SetRune(x, y int, ch rune) error {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return errors.New("invalid point")
	}
	termbox.SetCell(x, y, ch, termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
	return nil
}

func (g *Gui) Rune(x, y int) (rune, error) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return 0, errors.New("invalid point")
	}
	c := termbox.CellBuffer()[y*g.maxX+x]
	return c.Ch, nil
}

func (g *Gui) SetView(name string, x0, y0, x1, y1 int) (*View, error) {
	if x0 >= x1 || y0 >= y1 {
		return nil, errors.New("invalid dimensions")
	}

	if v := g.View(name); v != nil {
		v.x0 = x0
		v.y0 = y0
		v.x1 = x1
		v.y1 = y1
		return v, nil
	}

	v := newView(name, x0, y0, x1, y1)
	v.bgColor, v.fgColor = g.BgColor, g.FgColor
	v.selBgColor, v.selFgColor = g.SelBgColor, g.SelFgColor
	g.views = append(g.views, v)
	return v, ErrorUnkView
}

func (g *Gui) View(name string) *View {
	for _, v := range g.views {
		if v.name == name {
			return v
		}
	}
	return nil
}

func (g *Gui) DeleteView(name string) error {
	for i, v := range g.views {
		if v.name == name {
			g.views = append(g.views[:i], g.views[i+1:]...)
			return nil
		}
	}
	return ErrorUnkView
}

func (g *Gui) SetCurrentView(name string) error {
	for _, v := range g.views {
		if v.name == name {
			g.currentView = v
			return nil
		}
	}
	return ErrorUnkView
}

func (g *Gui) SetKeybinding(viewname string, key interface{}, mod Modifier, cb KeybindingCB) error {
	var kb *keybinding

	switch k := key.(type) {
	case Key:
		kb = newKeybinding(viewname, k, 0, mod, cb)
	case rune:
		kb = newKeybinding(viewname, 0, k, mod, cb)
	default:
		return errors.New("unknown type")
	}
	g.keybindings = append(g.keybindings, kb)
	return nil
}

func (g *Gui) SetLayout(layout func(*Gui) error) {
	g.layout = layout
	g.currentView = nil
	g.views = nil
	go func() { g.events <- termbox.Event{Type: termbox.EventResize} }()
}

func (g *Gui) MainLoop() error {
	go func() {
		for {
			g.events <- termbox.PollEvent()
		}
	}()

	termbox.SetInputMode(termbox.InputAlt)

	if err := g.Flush(); err != nil {
		return err
	}
	for {
		ev := <-g.events
		if err := g.handleEvent(&ev); err != nil {
			return err
		}
		if err := g.consumeevents(); err != nil {
			return err
		}
		if err := g.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gui) consumeevents() error {
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

func (g *Gui) handleEvent(ev *termbox.Event) error {
	switch ev.Type {
	case termbox.EventKey:
		return g.onKey(ev)
	case termbox.EventError:
		return ev.Err
	default:
		return nil
	}
}

func (g *Gui) draw(v *View) error {
	if g.ShowCursor {
		if v := g.currentView; v != nil {
			maxX, maxY := v.Size()
			cx, cy := v.Cursor()
			if v.cx >= maxX {
				cx = maxX - 1
			}
			if v.cy >= maxY {
				cy = maxY - 1
			}
			if err := v.SetCursor(cx, cy); err != nil {
				return nil
			}
			termbox.SetCursor(v.x0+v.cx+1, v.y0+v.cy+1)
		}
	} else {
		termbox.HideCursor()
	}

	v.clearRunes()
	if err := v.draw(); err != nil {
		return err
	}
	return nil
}

func (g *Gui) Flush() error {
	if g.layout == nil {
		return errors.New("Null layout")
	}

	termbox.Clear(termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
	g.maxX, g.maxY = termbox.Size()
	if err := g.layout(g); err != nil {
		return err
	}
	for _, v := range g.views {
		if err := g.drawFrame(v); err != nil {
			return err
		}
		if err := g.draw(v); err != nil {
			return err
		}
	}
	if err := g.drawIntersections(); err != nil {
		return err
	}
	termbox.Flush()
	return nil

}

func (g *Gui) drawFrame(v *View) error {
	for x := v.x0 + 1; x < v.x1 && x < g.maxX; x++ {
		if x < 0 {
			continue
		}
		if v.y0 > -1 && v.y0 < g.maxY {
			if err := g.SetRune(x, v.y0, '─'); err != nil {
				return err
			}
		}
		if v.y1 > -1 && v.y1 < g.maxY {
			if err := g.SetRune(x, v.y1, '─'); err != nil {
				return err
			}
		}
	}
	for y := v.y0 + 1; y < v.y1 && y < g.maxY; y++ {
		if y < 0 {
			continue
		}
		if v.x0 > -1 && v.x0 < g.maxX {
			if err := g.SetRune(v.x0, y, '│'); err != nil {
				return err
			}
		}
		if v.x1 > -1 && v.x1 < g.maxX {
			if err := g.SetRune(v.x1, y, '│'); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Gui) drawIntersections() error {
	for _, v := range g.views {
		if ch, ok := g.intersectionRune(v.x0, v.y0); ok {
			if err := g.SetRune(v.x0, v.y0, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.intersectionRune(v.x0, v.y1); ok {
			if err := g.SetRune(v.x0, v.y1, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.intersectionRune(v.x1, v.y0); ok {
			if err := g.SetRune(v.x1, v.y0, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.intersectionRune(v.x1, v.y1); ok {
			if err := g.SetRune(v.x1, v.y1, ch); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Gui) intersectionRune(x, y int) (rune, bool) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return 0, false
	}

	chTop, _ := g.Rune(x, y-1)
	top := verticalRune(chTop)
	chBottom, _ := g.Rune(x, y+1)
	bottom := verticalRune(chBottom)
	chLeft, _ := g.Rune(x-1, y)
	left := horizontalRune(chLeft)
	chRight, _ := g.Rune(x+1, y)
	right := horizontalRune(chRight)

	var ch rune
	switch {
	case !top && bottom && !left && right:
		ch = '┌'
	case !top && bottom && left && !right:
		ch = '┐'
	case top && !bottom && !left && right:
		ch = '└'
	case top && !bottom && left && !right:
		ch = '┘'
	case top && bottom && left && right:
		ch = '┼'
	case top && bottom && !left && right:
		ch = '├'
	case top && bottom && left && !right:
		ch = '┤'
	case !top && bottom && left && right:
		ch = '┬'
	case top && !bottom && left && right:
		ch = '┴'
	default:
		return 0, false
	}
	return ch, true
}

func verticalRune(ch rune) bool {
	if ch == '│' || ch == '┼' || ch == '├' || ch == '┤' {
		return true
	}
	return false
}

func horizontalRune(ch rune) bool {
	if ch == '─' || ch == '┼' || ch == '┬' || ch == '┴' {
		return true
	}
	return false
}

func (g *Gui) onKey(ev *termbox.Event) error {
	for _, kb := range g.keybindings {
		if ev.Ch == kb.Ch && Key(ev.Key) == kb.Key && Modifier(ev.Mod) == kb.Mod &&
			(kb.ViewName == "" || (g.currentView != nil && kb.ViewName == g.currentView.name)) {
			if kb.CB == nil {
				return nil
			}
			return kb.CB(g, g.currentView)
		}
	}
	return nil
}

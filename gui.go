package gocui

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var ErrorQuit error = errors.New("quit")

type Gui struct {
	events           chan termbox.Event
	CurrentView      *View
	views            []*View
	keybindings      []*Keybinding
	Layout           func(*Gui) error
	Start            func(*Gui) error
	maxX, maxY       int
	BgColor, FgColor termbox.Attribute
	ShowCursor       bool
}

func NewGui() (g *Gui) {
	return &Gui{}
}

func (g *Gui) Init() (err error) {
	if err = termbox.Init(); err != nil {
		return err
	}
	g.events = make(chan termbox.Event, 20)
	g.maxX, g.maxY = termbox.Size()
	g.BgColor = termbox.ColorWhite
	g.FgColor = termbox.ColorBlack
	return nil
}

func (g *Gui) Close() {
	termbox.Close()
}

func (g *Gui) Size() (x, y int) {
	return g.maxX, g.maxY
}

func (g *Gui) SetRune(x, y int, ch rune) (err error) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return errors.New("invalid point")
	}
	termbox.SetCell(x, y, ch, g.FgColor, g.BgColor)
	return nil
}

func (g *Gui) GetRune(x, y int) (ch rune, err error) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return 0, errors.New("invalid point")
	}
	c := termbox.CellBuffer()[y*g.maxX+x]
	return c.Ch, nil
}

func (g *Gui) SetView(name string, x0, y0, x1, y1 int) (v *View, err error) {
	if x0 >= x1 || y0 >= y1 {
		return nil, errors.New("invalid dimensions")
	}

	if v := g.GetView(name); v != nil {
		v.X0 = x0
		v.Y0 = y0
		v.X1 = x1
		v.Y1 = y1
		return v, nil
	}

	v = NewView(name, x0, y0, x1, y1)
	g.views = append(g.views, v)
	return v, nil
}

func (g *Gui) GetView(name string) (v *View) {
	for _, v := range g.views {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (g *Gui) DeleteView(name string) (err error) {
	for i, v := range g.views {
		if v.Name == name {
			g.views = append(g.views[:i], g.views[i+1:]...)
			return nil
		}
	}
	return errors.New("unknown view")
}

func (g *Gui) SetCurrentView(name string) (err error) {
	for _, v := range g.views {
		if v.Name == name {
			g.CurrentView = v
			return nil
		}
	}
	return errors.New("unknown view")
}

func (g *Gui) SetKeybinding(viewname string, key interface{}, mod Modifier, cb KeybindingCB) (err error) {
	var kb *Keybinding

	switch k := key.(type) {
	case Key:
		kb = NewKeybinding(viewname, k, 0, mod, cb)
	case rune:
		kb = NewKeybinding(viewname, 0, k, mod, cb)
	default:
		return errors.New("unknown type")
	}
	g.keybindings = append(g.keybindings, kb)
	return nil
}

func (g *Gui) MainLoop() (err error) {
	go func() {
		for {
			g.events <- termbox.PollEvent()
		}
	}()

	termbox.SetInputMode(termbox.InputAlt)

	if err := g.resize(); err != nil {
		return err
	}
	if g.Start != nil {
		if err := g.Start(g); err != nil {
			return err
		}
	}
	if err := g.draw(); err != nil {
		return err
	}

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
		if err := g.drawView(v); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gui) drawView(v *View) (err error) {
	if g.ShowCursor && v == g.CurrentView {
		termbox.SetCursor(v.X0+v.CX+1, v.Y0+v.CY+1)
	}

	return nil
}

func (g *Gui) resize() (err error) {
	if g.Layout == nil {
		return errors.New("Null layout")
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	g.maxX, g.maxY = termbox.Size()
	if err := g.Layout(g); err != nil {
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
	for _, v := range g.views {
		for x := v.X0 + 1; x < v.X1 && x < g.maxX; x++ {
			if x < 0 {
				continue
			}
			if v.Y0 > -1 && v.Y0 < g.maxY {
				if err := g.SetRune(x, v.Y0, '─'); err != nil {
					return err
				}
			}
			if v.Y1 > -1 && v.Y1 < g.maxY {
				if err := g.SetRune(x, v.Y1, '─'); err != nil {
					return err
				}
			}
		}
		for y := v.Y0 + 1; y < v.Y1 && y < g.maxY; y++ {
			if y < 0 {
				continue
			}
			if v.X0 > -1 && v.X0 < g.maxX {
				if err := g.SetRune(v.X0, y, '│'); err != nil {
					return err
				}
			}
			if v.X1 > -1 && v.X1 < g.maxX {
				if err := g.SetRune(v.X1, y, '│'); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *Gui) drawIntersections() (err error) {
	for _, v := range g.views {
		if ch, ok := g.getIntersectionRune(v.X0, v.Y0); ok {
			if err := g.SetRune(v.X0, v.Y0, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.getIntersectionRune(v.X0, v.Y1); ok {
			if err := g.SetRune(v.X0, v.Y1, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.getIntersectionRune(v.X1, v.Y0); ok {
			if err := g.SetRune(v.X1, v.Y0, ch); err != nil {
				return err
			}
		}
		if ch, ok := g.getIntersectionRune(v.X1, v.Y1); ok {
			if err := g.SetRune(v.X1, v.Y1, ch); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Gui) getIntersectionRune(x, y int) (ch rune, ok bool) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return 0, false
	}

	chTop, _ := g.GetRune(x, y-1)
	top := verticalRune(chTop)
	chBottom, _ := g.GetRune(x, y+1)
	bottom := verticalRune(chBottom)
	chLeft, _ := g.GetRune(x-1, y)
	left := horizontalRune(chLeft)
	chRight, _ := g.GetRune(x+1, y)
	right := horizontalRune(chRight)

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

func (g *Gui) onKey(ev *termbox.Event) (err error) {
	for _, kb := range g.keybindings {
		if ev.Ch == kb.Ch && Key(ev.Key) == kb.Key && Modifier(ev.Mod) == kb.Mod &&
			(kb.ViewName == "" || (g.CurrentView != nil && kb.ViewName == g.CurrentView.Name)) {
			if err := kb.CB(g, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

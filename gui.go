// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var (
	// ErrorQuit is used to decide if the MainLoop finished succesfully.
	ErrorQuit error = errors.New("quit")

	// ErrorUnkView allows to assert if a View must be initialized.
	ErrorUnkView error = errors.New("unknown view")
)

// Gui represents the whole User Interface, including the views, layouts
// and keybindings.
type Gui struct {
	events      chan termbox.Event
	views       []*View
	currentView *View
	layout      func(*Gui) error
	keybindings []*keybinding
	maxX, maxY  int
	redrawBgFg  bool // hook for redraw the background and foreground

	// BgColor and FgColor allow to configure the background and foreground
	// colors of the GUI.
	BgColor, FgColor Attribute

	// SelBgColor and SelFgColor are used to configure the background and
	// foreground colors of the selected line, when it is highlighted.
	SelBgColor, SelFgColor Attribute

	// If ShowCursor is true then the cursor is enabled.
	ShowCursor bool
}

// NewGui returns a new Gui object.
func NewGui() *Gui {
	return &Gui{}
}

// Init initializes the library. This function must be called before
// any other functions.
func (g *Gui) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	g.events = make(chan termbox.Event, 20)
	g.maxX, g.maxY = termbox.Size()
	g.BgColor = ColorBlack
	g.FgColor = ColorWhite
	g.redrawBgFg = true
	return nil
}

// Close finalizes the library. It should be called after a successful
// initialization and when gocui is not needed anymore.
func (g *Gui) Close() {
	termbox.Close()
}

// Size returns the terminal's size.
func (g *Gui) Size() (x, y int) {
	return g.maxX, g.maxY
}

// SetRune writes a rune at the given point, relative to the top-left
// corner of the terminal. It checks if the position is valid and applies
// the gui's colors.
func (g *Gui) SetRune(x, y int, ch rune) error {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return errors.New("invalid point")
	}
	termbox.SetCell(x, y, ch, termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
	return nil
}

// Rune returns the rune contained in the cell at the given position.
// It checks if the position is valid.
func (g *Gui) Rune(x, y int) (rune, error) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return ' ', errors.New("invalid point")
	}
	c := termbox.CellBuffer()[y*g.maxX+x]
	return c.Ch, nil
}

// SetView creates a new view with its top-left corner at (x0, y0)
// and the bottom-right one at (x1, y1). If a view with the same name
// already exists, its dimensions are updated; otherwise, the error
// ErrorUnkView is returned, which allows to assert if the View must
// be initialized. It checks if the position is valid.
func (g *Gui) SetView(name string, x0, y0, x1, y1 int) (*View, error) {
	if x0 >= x1 || y0 >= y1 {
		return nil, errors.New("invalid dimensions")
	}
	if name == "" {
		return nil, errors.New("invalid name")
	}

	if v := g.View(name); v != nil {
		// compare coordinates
		// if they are different - need to redraw the View
		if v.x0 == x0 && v.x1 == x1 && v.y0 == y0 && v.y1 == y1 {
			return v, nil
		}

		vxy := &View{
			x0: x0,
			y0: y0,
			x1: x1,
			y1: y1,
		}
		// compare old and new coordinates with the overlapping of other Views
		for _, r := range g.views {
			if v.name == r.name {
				continue
			}

			oldCross := g.checkOwerlap(vxy, v)
			newCross := g.checkOwerlap(vxy, r)
			var x0Cross, x1Cross, y0Cross, y1Cross bool

			if oldCross && !newCross {
				r.redraw = true
				continue
			}

			if (r.x0 < v.x0 && v.x0 < x0) || (r.x0 > v.x0 && x0 > r.x0) {
				x0Cross = true
			}
			if (r.y0 < v.y0 && v.y0 < y0) || (r.y0 > v.y0 && y0 > r.y0) {
				y0Cross = true
			}
			if (r.x1 > v.x1 && v.x1 > x1) || (r.x1 < v.x1 && x1 < r.x1) {
				x1Cross = true
			}
			if (r.y1 > v.y1 && v.y1 > y1) || (r.y1 < v.y1 && y1 < r.y1) {
				y1Cross = true
			}

			if x0Cross || y0Cross || x1Cross || y1Cross {
				r.redraw = true
			}
		}
		// create View with old coordinates and set "delete" option;
		// this View will "close holes" in the background
		for i, z := range g.views {
			if z.name == v.name+"-delmask" {
				g.views = append(g.views[:i], g.views[i+1:]...)
			}
		}
		k := newView(v.name+"-delmask", v.x0-1, v.y0-1, v.x1+1, v.y1+1)
		k.Frame = false
		k.BgColor, k.FgColor = g.BgColor, g.FgColor
		k.deleted = true
		g.views = append(g.views, k)

		// accept the new coordinates
		v.x0 = x0
		v.y0 = y0
		v.x1 = x1
		v.y1 = y1
		v.redraw = true

		return v, nil
	}

	v := newView(name, x0, y0, x1, y1)
	v.BgColor, v.FgColor = g.BgColor, g.FgColor
	v.SelBgColor, v.SelFgColor = g.SelBgColor, g.SelFgColor
	v.redraw = true
	g.views = append(g.views, v)
	return v, ErrorUnkView
}

// View returns a pointer to the view with the given name, or nil if
// a view with that name does not exist.
func (g *Gui) View(name string) *View {
	for _, v := range g.views {
		if v.name == name {
			return v
		}
	}
	return nil
}

// DeleteView deletes a view by name.
func (g *Gui) DeleteView(name string) error {
	for _, v := range g.views {
		if v.name != name {
			continue
		}
		// clear and change deleted View properties
		v.Clear()
		v.x0 = v.x0 - 1 // for borders overlap
		v.y0 = v.y0 - 1
		v.x1 = v.x1 + 1
		v.y1 = v.y1 + 1
		v.Frame = false
		v.BgColor = g.BgColor
		v.FgColor = g.FgColor
		v.deleted = true

		// if the deleted View was drawn over the other Views (or some parts
		// of them), we must to redraw these Views
		dvxy := &View{
			x0: v.x0,
			y0: v.y0,
			x1: v.x1,
			y1: v.y1,
		}
		// compare the coordinates of the deleted View with existing
		for _, r := range g.views {
			if g.checkOwerlap(dvxy, r) {
				r.redraw = true
			}
		}
		return nil

	}
	return ErrorUnkView
}

// SetCurrentView gives the focus to a given view.
func (g *Gui) SetCurrentView(name string) error {
	for _, v := range g.views {
		if v.name == name {
			g.currentView = v
			v.redraw = true
			return nil
		}
	}
	return ErrorUnkView
}

// CurrentView returns the currently focused view, or nil if no view
// owns the focus.
func (g *Gui) CurrentView() *View {
	return g.currentView
}

// SetKeybinding creates a new keybinding. If viewname equals to ""
// (empty string) then the keybinding will apply to all views. key must
// be a rune or a Key.
func (g *Gui) SetKeybinding(viewname string, key interface{}, mod Modifier, h KeybindingHandler) error {
	var kb *keybinding

	switch k := key.(type) {
	case Key:
		kb = newKeybinding(viewname, k, 0, mod, h)
	case rune:
		kb = newKeybinding(viewname, 0, k, mod, h)
	default:
		return errors.New("unknown type")
	}
	g.keybindings = append(g.keybindings, kb)
	return nil
}

// SetLayout sets the current layout. A layout is a function that
// will be called everytime the gui is re-drawed, it must contain
// the base views and its initializations.
func (g *Gui) SetLayout(layout func(*Gui) error) {
	g.layout = layout
	g.currentView = nil
	g.views = nil
	go func() { g.events <- termbox.Event{Type: termbox.EventResize} }()
}

// MainLoop runs the main loop until an error is returned. A successful
// finish should return ErrorQuit.
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

// consumeevents handles the remaining events in the events pool.
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

// handleEvent handles an event, based on its type (key-press, error,
// etc.)
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

// Flush updates the gui, re-drawing frames and buffers.
func (g *Gui) Flush() error {
	if g.layout == nil {
		return errors.New("Null layout")
	}

	if g.redrawBgFg {
		termbox.Clear(termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
		for _, v := range g.views {
			v.redraw = true
		}
		g.redrawBgFg = false
	}

	// check console size change
	// and redraw all Views if the size changed
	mX, mY := termbox.Size()
	if mX != g.maxX || mY != g.maxY {
		termbox.Clear(termbox.Attribute(g.FgColor), termbox.Attribute(g.BgColor))
		for _, v := range g.views {
			v.redraw = true
		}
	}
	g.maxX, g.maxY = mX, mY

	// сreate pools for alternately draw
	var deletedViews []*View // must be redrawn first of all
	var normalViews []*View
	var topLevelViews []*View // must be redrawn after the normal Views

	if err := g.layout(g); err != nil {
		return err
	}
	for i, v := range g.views {
		switch {
		case v.deleted:
			deletedViews = append(deletedViews, v)
			g.views = append(g.views[:i], g.views[i+1:]...)
		case v.AlwaysOnTop:
			topLevelViews = append(topLevelViews, v)
		case v.redraw:
			normalViews = append(normalViews, v)
		}
	}
	drawList := append(append(deletedViews, normalViews...), topLevelViews...)

	for _, v := range drawList {
		if v.Frame {
			if err := g.drawFrame(v); err != nil {
				return err
			}
		}
		if err := g.draw(v); err != nil {
			return err
		}
		v.redraw = false
	}

	//for i, v := range g.views {
	//	if v.deleted {
	//		g.views = append(g.views[:i], g.views[i+1:]...)
	//	}
	//}

	if err := g.drawIntersections(); err != nil {
		return err
	}
	termbox.Flush()
	return nil

}

// drawFrame draws the horizontal and vertical edges of a view.
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

// draw manages the cursor and calls the draw function of a view.
func (g *Gui) draw(v *View) error {
	if g.ShowCursor {
		if v := g.currentView; v != nil {
			maxX, maxY := v.Size()
			cx, cy := v.cx, v.cy
			if v.cx >= maxX {
				cx = maxX - 1
			}
			if v.cy >= maxY {
				cy = maxY - 1
			}
			if err := v.SetCursor(cx, cy); err != nil {
				return err
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

// drawIntersections draws the corners of each view, based on the type
// of the edges that converge at these points.
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

// intersectionRune returns the correct intersection rune at a given
// point.
func (g *Gui) intersectionRune(x, y int) (rune, bool) {
	if x < 0 || y < 0 || x >= g.maxX || y >= g.maxY {
		return ' ', false
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
		return ' ', false
	}
	return ch, true
}

// verticalRune returns if the given character is a vertical rune.
func verticalRune(ch rune) bool {
	if ch == '│' || ch == '┼' || ch == '├' || ch == '┤' {
		return true
	}
	return false
}

// verticalRune returns if the given character is a horizontal rune.
func horizontalRune(ch rune) bool {
	if ch == '─' || ch == '┼' || ch == '┬' || ch == '┴' {
		return true
	}
	return false
}

// onKey manages key-press events. A keybinding handler is called when
// a key-press event satisfies a configured keybinding. Furthermore,
// currentView's internal buffer is modified if currentView.Editable is true.
func (g *Gui) onKey(ev *termbox.Event) error {
	if g.currentView != nil && g.currentView.Editable {
		if err := g.handleEdit(g.currentView, ev); err != nil {
			return err
		}
	}
	for _, kb := range g.keybindings {
		if kb.h != nil && ev.Ch == kb.ch && Key(ev.Key) == kb.key && Modifier(ev.Mod) == kb.mod &&
			(kb.viewName == "" || (g.currentView != nil && kb.viewName == g.currentView.name)) {
			return kb.h(g, g.currentView)
		}
	}
	return nil
}

// handleEdit manages the edition mode.
func (g *Gui) handleEdit(v *View, ev *termbox.Event) error {
	switch {
	case ev.Ch != 0 && ev.Mod == 0:
		return v.editWrite(ev.Ch)
	case ev.Key == termbox.KeySpace:
		return v.editWrite(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		return v.editDelete(true)
	case ev.Key == termbox.KeyDelete:
		return v.editDelete(false)
	case ev.Key == termbox.KeyInsert:
		v.overwrite = !v.overwrite
	case ev.Key == termbox.KeyEnter:
		return v.editLine()
	}
	return nil
}

// checkOwerlap checks whether overlapping Views or not
func (g *Gui) checkOwerlap(v0, v1 *View) bool {
	var xCross, yCross bool
	if (v0.x0 < v1.x0 && v0.x1 > v1.x0) || (v0.x0 > v1.x0 && v0.x0 < v1.x1) {
		xCross = true
	}
	if (v0.y0 < v1.y0 && v0.y1 > v1.y0) || (v0.y0 > v1.y0 && v0.y0 < v1.y1) {
		yCross = true
	}
	if xCross && yCross {
		return true
	}
	return false
}

package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

var (
	viewArr = []string{"v1", "v2", "v3", "v4"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("v1")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX/2-1, maxY/2-1, gocui.RIGHT); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "v1"
		v.Autoscroll = true
		fmt.Fprintln(v, "View with default frame color")
		fmt.Fprintln(v, "It's connected to v2 with overlay RIGHT.\n")
		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", maxX/2-1, 0, maxX-1, maxY/2-1, gocui.LEFT); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "v2"
		v.Wrap = true
		v.FrameColor = gocui.ColorMagenta
		v.FrameRunes = []rune{'═', '│'}
		fmt.Fprintln(v, "View with minimum frame customization and colored frame.")
		fmt.Fprintln(v, "It's connected to v1 with overlay LEFT.\n")
		fmt.Fprintln(v, "\033[35;1mInstructions:\033[0m")
		fmt.Fprintln(v, "Press TAB to change current view")
		fmt.Fprintln(v, "Press Ctrl+O to toggle gocui.SupportOverlap\n")
		fmt.Fprintln(v, "\033[32;2mSelected frame is highlighted with green color\033[0m")
	}
	if v, err := g.SetView("v3", 0, maxY/2, maxX/2-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "v3"
		v.Wrap = true
		v.Autoscroll = true
		v.FrameColor = gocui.ColorCyan
		v.TitleColor = gocui.ColorCyan
		v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝'}
		fmt.Fprintln(v, "View with basic frame customization and colored frame and title")
		fmt.Fprintln(v, "It's not connected to any view.")
	}
	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1, gocui.LEFT); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "v4"
		v.Subtitle = "(editable)"
		v.Editable = true
		v.TitleColor = gocui.ColorYellow
		v.FrameColor = gocui.ColorRed
		v.FrameRunes = []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}
		fmt.Fprintln(v, "View with fully customized frame and colored title differently.")
		fmt.Fprintln(v, "It's connected to v3 with overlay LEFT.\n")
		v.SetCursor(0, 3)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func toggleOverlap(g *gocui.Gui, v *gocui.View) error {
	g.SupportOverlaps = !g.SupportOverlaps
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	g.SelFrameColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, toggleOverlap); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

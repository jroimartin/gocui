package main

import (
	"log"

	"github.com/jesseduffield/gocui"
)

func main() {

	g, err := gocui.NewGui(gocui.OutputNormal, true, gocui.NORMAL, false)
	g.Cursor = true
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", nil, gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	if v, err := g.SetView("prompt", 0, 0, 10, 3, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = "prompt"
		v.Editable = true
		v.Wrap = true
		_, err := g.SetCurrentView(v.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

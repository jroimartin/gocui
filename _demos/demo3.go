package main

import (
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	var err error

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	if _, err := g.AddView("side", 0, 0, 0.3, 0.9); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("main", 0.3, 0, 0.7, 0.9); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("cmdline", 0, 0.9, 0.5, 0.1); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}

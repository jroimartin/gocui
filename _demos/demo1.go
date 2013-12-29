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

	maxX, maxY := g.Size()
	if _, err := g.AddView("v1", -1, -1, 10, 10); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v2", maxX-10, -1, maxX, 10); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v3", maxX/2-5, -1, maxX/2+5, 10); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v4", -1, maxY/2-5, 10, maxY/2+5); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v5", maxX-10, maxY/2-5, maxX, maxY/2+5); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v6", -1, maxY-10, 10, maxY); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v7", maxX-10, maxY-10, maxX, maxY); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("v8", maxX/2-5, maxY-10, maxX/2+5, maxY); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}

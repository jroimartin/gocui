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
	if _, err := g.AddView("side", -1, -1, 31, float32(maxY-4)); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("main", 30, -1, float32(maxX-30), float32(maxY-4)); err != nil {
		log.Panicln(err)
	}
	if _, err := g.AddView("cmdline", -1, float32(maxY-5), float32(maxX+1), 5); err != nil {
		log.Panicln(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrorQuit {
		log.Panicln(err)
	}
}

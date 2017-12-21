package main

import (
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/jroimartin/gocui/table"
)

func fmtTime(v interface{}) string {
	t := v.(time.Time)
	return t.Format("2006-01-02 15:04:05")
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("new gocui: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Fatalf("keybindings: %v", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("table", 1, 1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		t := table.New().SetWidth(maxX)

		t.AddCol("Created").SetFormatFn(fmtTime)
		t.AddCol("Name").SetWidthPerc(100)
		t.AddCol("Age").AlignRight()
		t.AddCol("City").SetWidthPerc(50)

		t.AddRow(time.Now(), "Peter", 23, "Chicago")
		t.AddRow(time.Now(), "Sara", 15, "San Francisco")
		t.AddRow(time.Now(), "Sara", 45, "New York")
		t.AddRow(time.Now(), "John", 23, "Newark")
		t.AddRow(time.Now(), "Ariana", 34, "Los Angeles")

		t.SortAsc("Name").SortDesc("Age").Sort().Format().Fprint(v)
	}
	return nil
}

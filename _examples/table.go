// Copyright 2017 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type Column struct {
	Title string
	Size  float32
}

type Table struct {
	name          string
	Left, Top     int
	Right, Bottom int
	Columns       []Column
	Data          [][]string
}

func NewTable(name string, left, top, right, bottom int) *Table {
	return &Table{
		name:   name,
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}
}

func (t *Table) Layout(g *gocui.Gui) error {
	view, err := g.SetView(t.name, t.Left, t.Top, t.Right, t.Bottom, 0)
	if err != nil && !gocui.IsUnknownView(err) {
		return err
	}

	width, height := view.Size()
	hOffset := 0
	for cid, column := range t.Columns {
		size := int(float32(width) * column.Size)

		view.SetWritePos(hOffset, 0)
		view.WriteString(column.Title)

		for rid := 0; rid < height; rid++ {
			if rid < len(t.Data[cid]) {
				view.SetWritePos(hOffset, rid+1)
				view.WriteString(t.Data[cid][rid])
			}
			view.SetWritePos(hOffset+size-3, rid)
			view.WriteRunes([]rune{'│'})
		}

		hOffset += size
	}

	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	table := NewTable("t", 1, 2, 80, 10)
	table.Columns = []Column{
		{"Column1", 0.25},
		{"Column2", 0.25},
		{"Column3", 0.25},
		{"Column4", 0.25},
	}
	table.Data = [][]string{
		{"00", "01", "02", "03"},
		{"10", "11", "12", "13"},
		{"20", "21", "22", "23"},
		{"30", "31", "32", "33"},
	}
	g.SetManager(table)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

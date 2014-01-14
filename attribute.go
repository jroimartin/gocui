// Copyright 2014 The gocui Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import "github.com/nsf/termbox-go"

type Attribute termbox.Attribute

const (
	ColorDefault Attribute = Attribute(termbox.ColorDefault)
	ColorBlack             = Attribute(termbox.ColorBlack)
	ColorRed               = Attribute(termbox.ColorRed)
	ColorGreen             = Attribute(termbox.ColorGreen)
	ColorYellow            = Attribute(termbox.ColorYellow)
	ColorBlue              = Attribute(termbox.ColorBlue)
	ColorMagenta           = Attribute(termbox.ColorMagenta)
	ColorCyan              = Attribute(termbox.ColorCyan)
	ColorWhite             = Attribute(termbox.ColorWhite)
)

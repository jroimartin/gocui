# GOCUI - Go Console User Interface

[![Go Reference](https://pkg.go.dev/badge/github.com/jroimartin/gocui.svg)](https://pkg.go.dev/github.com/jroimartin/gocui)

Minimalist Go package aimed at creating Console User Interfaces.

## 关于本fork的说明（About this fork）

Add CJK support, both edit.go and view.go. 
修正了中文显示功能，为此我修改了`edit.go`和`view.go`。
因为该补丁还没有被原作者采用，使用时需要用用本代码覆盖 `github.com/jroimartin/gocui` 的代码。
2018-8-31 增加 `（* view）ReadEditor()`，能够正确过滤输入的中文，删除自动添加的占位符。

## Features

* Minimalist API.
* Views (the "windows" in the GUI) implement the interface io.ReadWriter.
* Support for overlapping views.
* The GUI can be modified at runtime (concurrent-safe).
* Global and view-level keybindings.
* Mouse support.
* Colored text.
* Customizable edition mode.
* Easy to build reusable widgets, complex layouts...

## Installation

Execute:

```sh
go get github.com/jroimartin/gocui
```

## Documentation

Execute:

```sh
go doc github.com/jroimartin/gocui
```

Or visit [pkg.go.dev](https://pkg.go.dev/github.com/jroimartin/gocui) to read
it online.

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
```

## Screenshots

![r2cui](https://cloud.githubusercontent.com/assets/1223476/19418932/63645052-93ce-11e6-867c-da5e97e37237.png)

![_examples/demo.go](https://cloud.githubusercontent.com/assets/1223476/5992750/720b84f0-aa36-11e4-88ec-296fa3247b52.png)

![_examples/dynamic.go](https://cloud.githubusercontent.com/assets/1223476/5992751/76ad5cc2-aa36-11e4-8204-6a90269db827.png)

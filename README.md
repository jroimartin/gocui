# GOCUI - Go Console User Interface

[![github actions](https://github.com/awesome-gocui/gocui/actions/workflows/go.yml/badge.svg)](https://github.com/awesome-gocui/gocui/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/awesome-gocui/gocui)](https://goreportcard.com/report/github.com/awesome-gocui/gocui)
[![GoDoc](https://godoc.org/github.com/awesome-gocui/gocui?status.svg)](https://godoc.org/github.com/awesome-gocui/gocui)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/awesome-gocui/gocui.svg)

Minimalist Go package aimed at creating Console User Interfaces.
A community fork based on the amazing work of [jroimartin](https://github.com/jroimartin/gocui)
For v0 to v1 mirgration help read: [migrate-to-v1.md](migrate-to-v1.md)

## Features

- Minimalist API.
- Views (the "windows" in the GUI) implement the interface io.ReadWriter.
- Support for overlapping views.
- The GUI can be modified at runtime (concurrent-safe).
- Global and view-level keybindings.
- Mouse support.
- Colored text.
- Customizable editing mode.
- Easy to build reusable widgets, complex layouts...

## About fork

This fork has many improvements over the original work from [jroimartin](https://github.com/jroimartin/gocui).

- Written ontop of TCell
- Better wide character support
- Support for 1 Line height views
- Support for running in docker container
- Better cursor handling
- Customize frame colors
- Improved code comments and quality
- Many small improvements
- Change Visibility of views
- Requires Go 1.13 or newer

For information about this org see: [awesome-gocui/about](https://github.com/awesome-gocui/about).

## Installation

Execute:

```
$ go get github.com/awesome-gocui/gocui
```

## Documentation

Execute:

```
$ go doc github.com/awesome-gocui/gocui
```

Or visit [godoc.org](https://godoc.org/github.com/awesome-gocui/gocui) to read it
online.

## Example

See the [\_example](./_example/) folder for more examples

```go
package main

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}

		if _, err := g.SetCurrentView("hello"); err != nil {
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

## Testing example

You can write simple tests for `gocui` which let you simulate keyboard and then validate the output drawn to the screen.

1. Create an instance of `gui` with `OutputSimulator` set as the mode `g, err := NewGui(OutputSimulator, true)`
2. Call `GetTestingScreen` to get a `testingScreen` instance. 
3. On this you can use `SendKey` to simulate input and `GetViewContent` to evaluate what is drawn.

> Warning: Timing plays a part here, key bindings don't fire synchronously and drawing isn't instant. Here we used `time.After` to pause, [`gomega`'s asynchronous assertions are likely a better alternative for more complex tests](https://onsi.github.io/gomega/#making-asynchronous-assertions).

Here is a simple example showing how this can be used to validate what a view shows and that a key binding is handled correctly:

```golang
func TestTestingScreenReturnsCorrectContent(t *testing.T) {
	// Track what happened in the view, we'll assert on these
	didCallCTRLC := false
	expectedViewContent := "Hello world!"
	viewName := "testView1"

	// Create a view specifying the "OutputSimulator" mode
	g, err := NewGui(OutputSimulator, true)
	if err != nil {
		log.Panicln(err)
	}
	g.SetManagerFunc(func(g *Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView(viewName, maxX/2-7, maxY/2, maxX/2+7, maxY/2+2, 0); err != nil {
			if !errors.Is(err, ErrUnknownView) {
				return err
			}

			if _, err := g.SetCurrentView(viewName); err != nil {
				return err
			}

			// Have the view draw "Hello world!"
			fmt.Fprintln(v, expectedViewContent)
		}

		return nil
	})

	// Create a key binding which sets "didCallCTRLC" when triggered
	exampleBindingToTest := func(g *Gui, v *View) error {
		didCallCTRLC = true
		return nil
	}
	if err := g.SetKeybinding("", KeyCtrlC, ModNone, exampleBindingToTest); err != nil {
		log.Panicln(err)
	}

	// Create a test screen and start gocui
	testingScreen := g.GetTestingScreen()
	cleanup := testingScreen.StartGui()
	defer cleanup()

	// Send a key to gocui
	testingScreen.SendKey(KeyCtrlC)

	// Wait for key to be processed
	<-time.After(time.Millisecond * 50)

	// Test that the keybinding fired and set "didCallCTRLC" to true
	if !didCallCTRLC {
		t.Error("Expect the simulator to invoke the key handler for CTRLC")
	}

	// Get the content from the testing screen
	actualContent, err := testingScreen.GetViewContent(viewName)
	if err != nil {
		t.Error(err)
	}

	// Test that it contains the "Hello World!" we thought the view should draw
	if strings.TrimSpace(actualContent) != expectedViewContent {
		t.Error(fmt.Printf("Expected view content to be: %q got: %q", expectedViewContent, actualContent))
	}
}

```

> Note: Under the covers this is using the `tcell` [`SimulationScreen`](https://github.com/gdamore/tcell/blob/master/simulation.go). 

## Screenshots

![r2cui](https://cloud.githubusercontent.com/assets/1223476/19418932/63645052-93ce-11e6-867c-da5e97e37237.png)

![_examples/demo.go](https://cloud.githubusercontent.com/assets/1223476/5992750/720b84f0-aa36-11e4-88ec-296fa3247b52.png)

![_examples/dynamic.go](https://cloud.githubusercontent.com/assets/1223476/5992751/76ad5cc2-aa36-11e4-8204-6a90269db827.png)

## Projects using gocui

- [komanda-cli](https://github.com/mephux/komanda-cli): IRC Client For Developers.
- [vuls](https://github.com/future-architect/vuls): Agentless vulnerability scanner for Linux/FreeBSD.
- [wuzz](https://github.com/asciimoo/wuzz): Interactive cli tool for HTTP inspection.
- [httplab](https://github.com/gchaincl/httplab): Interactive web server.
- [domainr](https://github.com/MichaelThessel/domainr): Tool that checks the availability of domains based on keywords.
- [gotime](https://github.com/nanohard/gotime): Time tracker for projects and tasks.
- [claws](https://github.com/thehowl/claws): Interactive command line client for testing websockets.
- [terminews](http://github.com/antavelos/terminews): Terminal based RSS reader.
- [diagram](https://github.com/esimov/diagram): Tool to convert ascii arts into hand drawn diagrams.
- [pody](https://github.com/JulienBreux/pody): CLI app to manage Pods in a Kubernetes cluster.
- [kubexp](https://github.com/alitari/kubexp): Kubernetes client.
- [kcli](https://github.com/cswank/kcli): Tool for inspecting kafka topics/partitions/messages.
- [fac](https://github.com/mkchoi212/fac): git merge conflict resolver
- [jsonui](https://github.com/gulyasm/jsonui): Interactive JSON explorer for your terminal.
- [cointop](https://github.com/miguelmota/cointop): Interactive terminal based UI application for tracking cryptocurrencies.
- [lazygit](https://github.com/jesseduffield/lazygit): simple terminal UI for git commands.
- [lazydocker](https://github.com/jesseduffield/lazydocker): The lazier way to manage everything docker.

Note: if your project is not listed here, let us know! :)

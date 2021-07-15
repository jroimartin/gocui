package gocui

import "time"

func (g *Gui) loaderTick() {
	go func() {
		for range time.Tick(time.Millisecond * 50) {
			for _, view := range g.Views() {
				if view.HasLoader {
					g.userEvents <- userEvent{func(g *Gui) error { return nil }}
					break
				}
			}
		}
	}()
}

// Loader can show a loading animation
func Loader() cell {
	characters := "|/-\\"
	now := time.Now()
	nanos := now.UnixNano()
	index := nanos / 50000000 % int64(len(characters))
	str := characters[index : index+1]
	chr := []rune(str)[0]
	return cell{
		chr: chr,
	}
}

package table

type FormatFn func(interface{}) string

type Col struct {
	name         string
	hide         bool
	format       string
	formatFn     FormatFn
	align        Align
	width        int
	perc         float32
	minWidth     int
	minWidthPerc int
}

type Cols []*Col

func (c *Col) Hide() *Col {
	c.hide = true
	return c
}

func (c *Col) SetFormat(f string) *Col {
	c.format = f
	return c
}

func (c *Col) SetFormatFn(f FormatFn) *Col {
	c.formatFn = f
	return c
}

func (c *Col) AlignLeft() *Col {
	c.align = AlignLeft
	return c
}

func (c *Col) AlignRight() *Col {
	c.align = AlignRight
	return c
}

func (c *Col) AlignCenter() *Col {
	c.align = AlignCenter
	return c
}

func (c *Col) SetWidth(w int) *Col {
	c.minWidth = w
	return c
}

func (c *Col) SetWidthPerc(w int) *Col {
	c.minWidthPerc = w
	return c
}

func (c Cols) Index(n string) int {
	for i := range c {
		if c[i].name == n {
			return i
		}
	}
	return -1
}

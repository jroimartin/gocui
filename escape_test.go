package gocui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOne(t *testing.T) {
	var ei *escapeInterpreter

	ei = newEscapeInterpreter(OutputNormal)
	isEscape, err := ei.parseOne('a')
	assert.Equal(t, false, isEscape)
	assert.NoError(t, err)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, []rune{'\x1b', '[', '0', 'K'})
	_, ok := ei.instruction.(eraseInLineFromCursor)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, []rune{'\x1b', '[', 'K'})
	_, ok = ei.instruction.(eraseInLineFromCursor)
	assert.Equal(t, true, ok)

	ei = newEscapeInterpreter(OutputNormal)
	parseEscRunes(t, ei, []rune{'\x1b', '[', '1', 'K'})
	_, ok = ei.instruction.(noInstruction)
	assert.Equal(t, true, ok)

}

func TestParseOneColours(t *testing.T) {
	scenarios := []struct {
		outputMode OutputMode
		runes      []rune
		expectedFg Attribute
		expectedBg Attribute
	}{
		{OutputNormal, []rune{'\x1b', '[', '3', '0', 'm'}, ColorBlack, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '1', 'm'}, ColorRed, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '2', 'm'}, ColorGreen, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '3', 'm'}, ColorYellow, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '4', 'm'}, ColorBlue, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '5', 'm'}, ColorMagenta, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '6', 'm'}, ColorCyan, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '3', '7', 'm'}, ColorWhite, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '4', '0', 'm'}, ColorDefault, ColorBlack},
		{OutputNormal, []rune{'\x1b', '[', '4', '1', 'm'}, ColorDefault, ColorRed},
		{OutputNormal, []rune{'\x1b', '[', '4', '2', 'm'}, ColorDefault, ColorGreen},
		{OutputNormal, []rune{'\x1b', '[', '4', '3', 'm'}, ColorDefault, ColorYellow},
		{OutputNormal, []rune{'\x1b', '[', '4', '4', 'm'}, ColorDefault, ColorBlue},
		{OutputNormal, []rune{'\x1b', '[', '4', '5', 'm'}, ColorDefault, ColorMagenta},
		{OutputNormal, []rune{'\x1b', '[', '4', '6', 'm'}, ColorDefault, ColorCyan},
		{OutputNormal, []rune{'\x1b', '[', '4', '7', 'm'}, ColorDefault, ColorWhite},
		{OutputNormal, []rune{'\x1b', '[', '4', '7', ';', '3', '1', 'm'}, ColorRed, ColorWhite},
	}

	for _, scenario := range scenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		parseEscRunes(t, ei, scenario.runes)
		assert.Equal(t, scenario.expectedFg, ei.curFgColor)
		assert.Equal(t, scenario.expectedBg, ei.curBgColor)
	}

	// resetting colours
	scenarios = []struct {
		outputMode OutputMode
		runes      []rune
		expectedFg Attribute
		expectedBg Attribute
	}{
		{OutputNormal, []rune{'\x1b', '[', '3', '9', 'm'}, ColorDefault, ColorRed},
		{OutputNormal, []rune{'\x1b', '[', '4', '9', 'm'}, ColorRed, ColorDefault},
		{OutputNormal, []rune{'\x1b', '[', '0', 'm'}, ColorDefault, ColorDefault},
	}

	for _, scenario := range scenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		ei.curFgColor = ColorRed
		ei.curBgColor = ColorRed
		parseEscRunes(t, ei, scenario.runes)
		assert.Equal(t, scenario.expectedFg, ei.curFgColor)
		assert.Equal(t, scenario.expectedBg, ei.curBgColor)
	}

	// setting attributes
	attrScenarios := []struct {
		outputMode   OutputMode
		runes        []rune
		expectedAttr Attribute
	}{
		{OutputNormal, []rune{'\x1b', '[', '1', 'm'}, AttrBold},
		{OutputNormal, []rune{'\x1b', '[', '2', 'm'}, AttrDim},
		{OutputNormal, []rune{'\x1b', '[', '3', 'm'}, AttrItalic},
		{OutputNormal, []rune{'\x1b', '[', '4', 'm'}, AttrUnderline},
		{OutputNormal, []rune{'\x1b', '[', '5', 'm'}, AttrBlink},
		{OutputNormal, []rune{'\x1b', '[', '7', 'm'}, AttrReverse},
		{OutputNormal, []rune{'\x1b', '[', '9', 'm'}, AttrStrikeThrough},
	}

	for _, scenario := range attrScenarios {
		ei := newEscapeInterpreter(scenario.outputMode)
		parseEscRunes(t, ei, scenario.runes)
		isBold := ei.curFgColor&scenario.expectedAttr == scenario.expectedAttr
		assert.Equal(t, true, isBold)
	}
}

func parseEscRunes(t *testing.T, ei *escapeInterpreter, runes []rune) {
	for _, r := range runes {
		isEscape, err := ei.parseOne(r)
		assert.Equal(t, true, isEscape)
		assert.NoError(t, err)
	}
}

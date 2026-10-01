// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocui

import (
	"fmt"
	"testing"
)

func parseString(ei *escapeInterpreter, str string) error {
	for _, r := range str {
		if _, err := ei.parseOne(r); err != nil {
			return err
		}
	}
	return nil
}

func TestExtendedANSIColorsNormal(t *testing.T) {
	tests := []struct {
		seq    string
		wantFg Attribute
		wantBg Attribute
	}{
		// Standard colors
		{"\x1b[30m", ColorBlack, ColorDefault},
		{"\x1b[31m", ColorRed, ColorDefault},
		{"\x1b[37m", ColorWhite, ColorDefault},
		{"\x1b[40m", ColorDefault, ColorBlack},
		{"\x1b[41m", ColorDefault, ColorRed},
		{"\x1b[47m", ColorDefault, ColorWhite},
		// Extended high-intensity foreground colors
		{"\x1b[90m", ColorDarkGray, ColorDefault},
		{"\x1b[91m", ColorLightRed, ColorDefault},
		{"\x1b[92m", ColorLightGreen, ColorDefault},
		{"\x1b[93m", ColorLightYellow, ColorDefault},
		{"\x1b[94m", ColorLightBlue, ColorDefault},
		{"\x1b[95m", ColorLightMagenta, ColorDefault},
		{"\x1b[96m", ColorLightCyan, ColorDefault},
		{"\x1b[97m", ColorLightGray, ColorDefault},
		// Extended high-intensity background colors
		{"\x1b[100m", ColorDefault, ColorDarkGray},
		{"\x1b[101m", ColorDefault, ColorLightRed},
		{"\x1b[102m", ColorDefault, ColorLightGreen},
		{"\x1b[103m", ColorDefault, ColorLightYellow},
		{"\x1b[104m", ColorDefault, ColorLightBlue},
		{"\x1b[105m", ColorDefault, ColorLightMagenta},
		{"\x1b[106m", ColorDefault, ColorLightCyan},
		{"\x1b[107m", ColorDefault, ColorLightGray},
		// Issue #241 reproducer: [0;91m
		{"\x1b[0;91m", ColorLightRed, ColorDefault},
		{"\x1b[0;101m", ColorDefault, ColorLightRed},
		// Combined bright fg and bright bg
		{"\x1b[91;100m", ColorLightRed, ColorDarkGray},
		{"\x1b[0;92;104m", ColorLightGreen, ColorLightBlue},
		// Default resets
		{"\x1b[39m", ColorDefault, ColorDefault},
		{"\x1b[49m", ColorDefault, ColorDefault},
	}

	for _, tc := range tests {
		ei := newEscapeInterpreter(OutputNormal)
		if err := parseString(ei, tc.seq); err != nil {
			t.Errorf("parseString(%q) failed: %v", tc.seq, err)
			continue
		}
		if ei.curFgColor != tc.wantFg {
			t.Errorf("seq %q: curFgColor = %v, want %v", tc.seq, ei.curFgColor, tc.wantFg)
		}
		if ei.curBgColor != tc.wantBg {
			t.Errorf("seq %q: curBgColor = %v, want %v", tc.seq, ei.curBgColor, tc.wantBg)
		}
	}
}

func TestExtendedANSIColors256(t *testing.T) {
	tests := []struct {
		seq    string
		wantFg Attribute
		wantBg Attribute
	}{
		// Extended high-intensity foreground colors in Output256 mode
		{"\x1b[90m", ColorDarkGray, ColorDefault},
		{"\x1b[91m", ColorLightRed, ColorDefault},
		{"\x1b[97m", ColorLightGray, ColorDefault},
		// Extended high-intensity background colors in Output256 mode
		{"\x1b[100m", ColorDefault, ColorDarkGray},
		{"\x1b[101m", ColorDefault, ColorLightRed},
		{"\x1b[107m", ColorDefault, ColorLightGray},
		// Issue #241 reproducer: [0;91m
		{"\x1b[0;91m", ColorLightRed, ColorDefault},
		{"\x1b[0;101m", ColorDefault, ColorLightRed},
		// Combined
		{"\x1b[91;100m", ColorLightRed, ColorDarkGray},
		{"\x1b[0;92;104m", ColorLightGreen, ColorLightBlue},
	}

	for _, tc := range tests {
		ei := newEscapeInterpreter(Output256)
		if err := parseString(ei, tc.seq); err != nil {
			t.Errorf("parseString(%q) failed: %v", tc.seq, err)
			continue
		}
		if ei.curFgColor != tc.wantFg {
			t.Errorf("seq %q in Output256: curFgColor = %v, want %v", tc.seq, ei.curFgColor, tc.wantFg)
		}
		if ei.curBgColor != tc.wantBg {
			t.Errorf("seq %q in Output256: curBgColor = %v, want %v", tc.seq, ei.curBgColor, tc.wantBg)
		}
	}
}

func TestViewWriteExtendedANSI(t *testing.T) {
	modes := []OutputMode{OutputNormal, Output256}

	for _, mode := range modes {
		t.Run(fmt.Sprintf("mode_%d", mode), func(t *testing.T) {
			v := newView("test", 0, 0, 20, 5, mode)
			input := "\x1b[0;91mHello\x1b[0m"
			n, err := v.Write([]byte(input))
			if err != nil {
				t.Fatalf("v.Write failed: %v", err)
			}
			if n != len(input) {
				t.Fatalf("v.Write returned %d, expected %d", n, len(input))
			}

			if len(v.lines) != 1 {
				t.Fatalf("expected 1 line, got %d", len(v.lines))
			}

			line := v.lines[0]
			if len(line) != 5 {
				t.Fatalf("expected 5 cells, got %d", len(line))
			}

			expectedText := "Hello"
			for i, c := range line {
				if c.chr != rune(expectedText[i]) {
					t.Errorf("cell %d: expected char %c, got %c", i, expectedText[i], c.chr)
				}
				if c.fgColor != ColorLightRed {
					t.Errorf("cell %d: expected fgColor %v (ColorLightRed), got %v", i, ColorLightRed, c.fgColor)
				}
				if c.bgColor != ColorDefault {
					t.Errorf("cell %d: expected bgColor %v, got %v", i, ColorDefault, c.bgColor)
				}
			}
		})
	}
}

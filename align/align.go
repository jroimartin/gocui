package align

import (
	"fmt"
	"strings"
)

func AlignLeft(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", n-len(s)))
}

func AlignRight(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", strings.Repeat(" ", n-len(s)), s)
}

func AlignCenter(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	pad := (n - len(s)) / 2
	lpad := pad
	rpad := n - len(s) - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), s, strings.Repeat(" ", rpad))
}

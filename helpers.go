package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func ToSnakeCase(s string) string {
	var b strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func progressBar(current, total, width int) string {
	if total <= 0 {
		return "[--------------------] 0%"
	}
	pct := current * 100 / total
	filled := current * width / total
	if filled > width {
		filled = width
	}
	return fmt.Sprintf("[%s%s] %3d%% (%d/%d)",
		strings.Repeat("#", filled),
		strings.Repeat("-", width-filled),
		pct, current, total)
}

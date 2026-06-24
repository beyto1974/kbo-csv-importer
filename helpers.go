package main

import (
	"bytes"
	"fmt"
	"io"
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

func printProgressBar(verbose bool, rowNum int, total int, lastShown int, newline bool) {
	if !verbose {
		return
	}

	p := ((rowNum + 1) * 100) / total
	if p != lastShown && (p == 100 || p%2 == 0) {
		fmt.Printf("\r  Progress: %s", progressBar(rowNum+1, total, 20))
		lastShown = p
	}

	if newline {
		fmt.Println()
	}
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	buf := make([]byte, 32*1024)
	nLines := 0

	for {
		n, err := f.Read(buf)
		nLines += bytes.Count(buf[:n], []byte{'\n'})

		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
	}

	return nLines, nil
}

// Package diff provides text diffing utilities.
package diff

import (
	"fmt"
	"strings"
)

// ANSI color codes.
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

// Op represents a diff operation type.
type Op int

const (
	OpEqual Op = iota
	OpDelete
	OpInsert
)

// Line represents a single line in a diff.
type Line struct {
	Op   Op
	Text string
}

// Compute computes a line-by-line diff using LCS algorithm.
func Compute(oldLines, newLines []string) []Line {
	m, n := len(oldLines), len(newLines)
	lcs := make([][]int, m+1)
	for i := range lcs {
		lcs[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				lcs[i][j] = lcs[i-1][j-1] + 1
			} else {
				lcs[i][j] = max(lcs[i-1][j], lcs[i][j-1])
			}
		}
	}

	// Backtrack to build diff
	var result []Line
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			result = append(result, Line{OpEqual, oldLines[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || lcs[i][j-1] >= lcs[i-1][j]) {
			result = append(result, Line{OpInsert, newLines[j-1]})
			j--
		} else {
			result = append(result, Line{OpDelete, oldLines[i-1]})
			i--
		}
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// SplitLines splits a string into lines, removing trailing newline.
func SplitLines(s string) []string {
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// PrintUnified prints a unified diff with color.
func PrintUnified(oldData, newData string) {
	oldLines := SplitLines(oldData)
	newLines := SplitLines(newData)

	diff := Compute(oldLines, newLines)

	const contextLines = 3
	lastPrinted := -1

	for i, d := range diff {
		if d.Op == OpEqual {
			continue
		}

		// Print context before
		start := i - contextLines
		if start < 0 {
			start = 0
		}
		if start <= lastPrinted {
			start = lastPrinted + 1
		}

		// Print separator if there's a gap
		if lastPrinted >= 0 && start > lastPrinted+1 {
			fmt.Printf("%s@@ ... @@%s\n", ColorCyan, ColorReset)
		}

		// Print leading context
		for j := start; j < i; j++ {
			if diff[j].Op == OpEqual {
				fmt.Printf("  %s\n", diff[j].Text)
				lastPrinted = j
			}
		}

		// Print the change
		printLine(d)
		lastPrinted = i

		// Print trailing context
		end := i + contextLines + 1
		if end > len(diff) {
			end = len(diff)
		}
		for j := i + 1; j < end; j++ {
			if diff[j].Op == OpEqual {
				fmt.Printf("  %s\n", diff[j].Text)
				lastPrinted = j
			} else {
				break
			}
		}
	}
}

func printLine(d Line) {
	switch d.Op {
	case OpDelete:
		fmt.Printf("%s- %s%s\n", ColorRed, d.Text, ColorReset)
	case OpInsert:
		fmt.Printf("%s+ %s%s\n", ColorGreen, d.Text, ColorReset)
	default:
		// OpEqual lines are handled separately in PrintUnified
	}
}

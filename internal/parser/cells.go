package parser

import "strings"

func SplitCells(notebookSQL string) []string {
	parts := strings.Split(notebookSQL, "\n\n")
	cells := make([]string, 0, len(parts))

	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			cells = append(cells, t)
		}
	}

	return cells
}

package notebook

import (
	"fmt"
	"strings"

	"github.com/kaustubha-chaturvedi/Rook/internal/parser"
)

func SQLFromFile(f File) string {
	if len(f.Cells) == 0 {
		return ""
	}

	parts := make([]string, 0, len(f.Cells))
	
	for _, c := range f.Cells {
		if q := strings.TrimSpace(c.Query); q != "" {
			parts = append(parts, q)
		}
	}

	return strings.Join(parts, "\n\n")
}

func FileFromSQL(sql string) File {
	parts := parser.SplitCells(sql)
	cells := make([]Cell, 0, len(parts))
	
	for i, q := range parts {
		cells = append(cells, Cell{
			ID:         fmt.Sprintf("cell-%d", i+1),
			Type:       "sql",
			Connection: "",
			Query:      q,
		})
	}
	
	return File{Version: CurrentVersion, Cells: cells}
}

package table

import (
	"fmt"
	"strings"
)

type render struct {
	tab       *Table
	lineLen   int
	separator string
}

func newRender(tab *Table) render {
	r := render{
		tab:     tab,
		lineLen: len(tab.lengths) + 1,
	}

	var separator strings.Builder

	for _, len := range tab.lengths {
		separator.WriteString("+")
		separator.WriteString(strings.Repeat("-", len+2))
		r.lineLen += len + 3
	}

	separator.WriteString("+\n")
	r.separator = separator.String()

	return r
}

func (r render) String() string {
	var sb strings.Builder

	sb.WriteString(r.separator)
	r.writeRow(&sb, r.tab.headers)

	sb.WriteString(r.separator)
	for _, row := range r.tab.rows {
		r.writeRow(&sb, row)
	}
	sb.WriteString(r.separator)

	return sb.String()
}

func (r render) writeRow(sb *strings.Builder, row []string) {
	for i, col := range row {
		sb.WriteString(fmt.Sprintf("| %s ", col))
		for j := len(col); j < r.tab.lengths[i]; j++ {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("|\n")
}

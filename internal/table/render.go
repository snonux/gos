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

	r.separator = tab.sprintf("%s+", separator.String()) + "\n"
	return r
}

func (r render) String() string {
	var sb strings.Builder

	sb.WriteString(r.separator)
	sb.WriteString(r.rowString(r.tab.headers))

	sb.WriteString(r.separator)
	for _, row := range r.tab.rows {
		sb.WriteString(r.rowString(row))
	}
	sb.WriteString(r.separator)

	return sb.String()
}

func (r render) rowString(row []string) string {
	var sb strings.Builder

	for i, col := range row {
		sb.WriteString(fmt.Sprintf("| %s ", col))
		for j := len(col); j < r.tab.lengths[i]; j++ {
			sb.WriteString(" ")
		}
	}

	return r.tab.sprintf("%s|", sb.String()) + "\n"
}

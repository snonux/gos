package table

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

type row []string

type Table struct {
	headers []string
	rows    []row
	lengths []int                                   // Max length of each col
	sprintf func(format string, args ...any) string // For colored output
	err     error
}

func New() *Table {
	return &Table{sprintf: fmt.Sprintf}
}

func (t *Table) WithColor(col *color.Color) *Table {
	t.sprintf = col.Sprintf
	return t
}

func (t *Table) Header(args ...any) *Table {
	t.headers = make([]string, 0, len(args))
	t.lengths = make([]int, 0, len(args))

	for _, arg := range args {
		strVal := val(arg)
		t.headers = append(t.headers, strVal)
		t.lengths = append(t.lengths, len(strVal))
	}

	return t
}

func (t *Table) Row(args ...any) *Table {
	if len(args) != len(t.headers) {
		t.err = fmt.Errorf("Table row (%v) not same length as table headers (%v)", args, t.headers)
	}
	if t.err != nil {
		return t
	}

	row := make(row, 0, len(args))
	for i, arg := range args {
		strVal := val(arg)
		row = append(row, strVal)
		if t.lengths[i] < len(row[i]) {
			t.lengths[i] = len(row[i])
		}
	}
	t.rows = append(t.rows, row)

	return t
}

func (t *Table) MustRender() {
	if err := t.Render(); err != nil {
		panic(err)
	}
}

func (t *Table) Render() error {
	if len(t.headers) == 0 {
		return fmt.Errorf("no headers")
	}
	if len(t.rows) == 0 {
		return fmt.Errorf("no rows")
	}
	if t.err != nil {
		return t.err
	}
	fmt.Print(newRender(t).String())
	return nil
}

func val(val any) string {
	switch v := val.(type) {
	case int:
		return strconv.Itoa(v)
	case float64:
		return fmt.Sprintf("%0.2f", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// func dataRow(sb *strings.Builder, descr1 string, val1 any, descr2 string, val2 any) {
// 	const format = "| %-21s | %-11s | %-21s | %-11s |"
// 	sb.WriteString(colour.SInfo2f(format, descr1, val(val1), descr2, val(val2)))
// 	sb.WriteString("\n")
// }

// func (t *Table) separator() {
// 	t.sb.WriteString(t.sep)
// 	t.sb.WriteString("\n")
// }

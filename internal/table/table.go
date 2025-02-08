package table

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/buger/goterm"
	"github.com/fatih/color"
)

type row []string

type formatFunc func(format string, args ...any) string

type Table struct {
	headers  []string
	rows     []row
	lengths  []int      // Max length of each col
	sheaderf formatFunc // For colored output
	sprintf  formatFunc // For colored output
	err      error
}

func New() *Table {
	return &Table{
		sprintf:  fmt.Sprintf,
		sheaderf: fmt.Sprintf,
	}
}

func (t *Table) WithColor(col *color.Color) *Table {
	return t.WithHeaderColor(col).WithBaseColor(col)
}

func (t *Table) WithHeaderColor(col *color.Color) *Table {
	t.sheaderf = col.Sprintf
	return t
}

func (t *Table) WithBaseColor(col *color.Color) *Table {
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
	t.addRow(vals(args...)...)
	return t
}

func (t *Table) TextBox(text string) *Table {
	maxLen := goterm.Width() - 4
	words := strings.Split(text, "\n")
	var result []string

	for _, line := range words {
		var currentLine string
		for _, word := range strings.Fields(line) {
			if len(currentLine)+len(word)+1 <= maxLen {
				if len(currentLine) > 0 {
					currentLine += " "
				}
				currentLine += word
			} else {
				if len(currentLine) > 0 {
					result = append(result, currentLine)
				}
				currentLine = word
			}
		}
		if len(currentLine) > 0 {
			result = append(result, currentLine)
		}
	}

	for _, line := range result {
		t.addRow(line)
	}

	return t
}

func (t *Table) addRow(row ...string) {
	if len(row) != len(t.headers) {
		t.err = fmt.Errorf("Table row (%v) not same length as table headers (%v)", row, t.headers)
	}
	if t.err != nil {
		return
	}

	for i, strVal := range row {
		if t.lengths[i] < len(strVal) {
			t.lengths[i] = len(strVal)
		}
	}
	t.rows = append(t.rows, row)
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

func vals(vals ...any) []string {
	strVals := make([]string, 0, len(vals))
	for _, v := range vals {
		strVals = append(strVals, val(v))
	}
	return strVals
}

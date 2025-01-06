package colour

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Infof = func(format string, args ...any) {
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
	Infoln = func(args ...any) { fmt.Println(args...) }

	info2Col = color.New(color.FgHiYellow, color.BgBlue)
	Info2f   = func(format string, args ...any) {
		info2Col.Printf(format, args...)
		fmt.Print("\n")
	}
	SInfo2f = info2Col.SprintfFunc()

	warnCol = color.New(color.FgHiWhite, color.BgRed)
	Warnln  = warnCol.PrintlnFunc()

	successCol = color.New(color.FgWhite, color.BgGreen)
	Successf   = func(format string, args ...any) {
		successCol.Printf(format, args...)
		fmt.Print("\n")

	}

	Ackf = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
)

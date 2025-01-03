package colour

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	// infoCol  = color.New(color.FgCyan, color.BgBlue, color.Bold)
	Infof    = func(format string, args ...any) { fmt.Printf(format, args...) }
	Infoln   = func(args ...any) { fmt.Println(args...) }
	info2Col = color.New(color.FgHiYellow, color.BgBlue)
	Info2f   = info2Col.PrintfFunc()
	SInfo2f  = info2Col.SprintfFunc()
	Ackf     = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
	warnCol  = color.New(color.FgHiWhite, color.BgRed)
	Warnf    = warnCol.PrintfFunc()
	Successf = color.New(color.FgWhite, color.BgGreen).PrintfFunc()
)

package colour

import "github.com/fatih/color"

var (
	infoCol  = color.New(color.FgCyan, color.BgBlue, color.Bold)
	Infof    = infoCol.PrintfFunc()
	Infoln   = infoCol.PrintlnFunc()
	info2Col = color.New(color.FgHiYellow, color.BgBlue)
	Info2f   = info2Col.PrintfFunc()
	SInfo2f  = info2Col.SprintfFunc()
	Ackf     = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
	Successf = color.New(color.FgWhite, color.BgGreen).PrintfFunc()
)

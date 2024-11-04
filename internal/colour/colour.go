package colour

import "github.com/fatih/color"

var (
	Info1f   = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	Info2f   = color.New(color.FgHiYellow, color.BgHiBlack, color.Bold).PrintfFunc()
	Sstatsf  = color.New(color.FgHiYellow, color.BgBlue).SprintfFunc()
	Ackf     = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
	Successf = color.New(color.FgWhite, color.BgGreen).PrintfFunc()
)

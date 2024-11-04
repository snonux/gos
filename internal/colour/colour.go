package colour

import "github.com/fatih/color"

var (
	// Printf function(s)
	Info1f   = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	Info2f   = color.New(color.FgHiYellow, color.BgHiBlack, color.Bold).PrintfFunc()
	Ackf     = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
	Successf = color.New(color.FgWhite, color.BgGreen).PrintfFunc()

	// Sprintf function(s)
	Sstatsf = color.New(color.FgHiYellow, color.BgBlue).SprintfFunc()
)

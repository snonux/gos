package colour

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Infofln = func(format string, args ...any) {
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
	Infoln = func(args ...any) { fmt.Println(args...) }

	Info2Col = color.New(color.FgHiYellow, color.BgBlue)
	Info2fln = func(format string, args ...any) {
		Info2Col.Printf(format, args...)
		fmt.Print("\n")
	}
	SInfo2f = Info2Col.SprintfFunc()

	warnCol = color.New(color.FgHiWhite, color.BgRed)
	Warnln  = warnCol.PrintlnFunc()

	successCol = color.New(color.FgWhite, color.BgGreen)
	Successfln = func(format string, args ...any) {
		successCol.Printf(format, args...)
		fmt.Print("\n")

	}

	Ackf = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
)

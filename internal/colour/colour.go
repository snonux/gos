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

	AttentionCol = color.New(color.FgHiYellow, color.BgBlue)
	warnCol      = color.New(color.FgHiWhite, color.BgRed)
	Warnln       = warnCol.PrintlnFunc()

	successCol = color.New(color.FgWhite, color.BgGreen)
	Successfln = func(format string, args ...any) {
		successCol.Printf(format, args...)
		fmt.Print("\n")

	}

	AckCol = color.New(color.FgBlack, color.BgHiYellow, color.Bold)
	Ackf   = AckCol.PrintfFunc()
)

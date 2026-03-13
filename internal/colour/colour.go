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

	errorCol = color.New(color.FgRed)
	Errorln  = errorCol.PrintlnFunc()

	successCol = color.New(color.FgWhite, color.BgGreen)
	Successfln = func(format string, args ...any) {
		if _, err := successCol.Printf(format, args...); err != nil {
			// Log the error but don't fail the operation since we've already printed the data
			Errorln("Error printing success message:", err)
		}
		fmt.Print("\n")

	}

	AckCol = color.New(color.FgBlack, color.BgHiYellow, color.Bold)
	Ackf   = AckCol.PrintfFunc()
)

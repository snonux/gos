package prompt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	Info1 = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	Info2 = color.New(color.FgHiYellow, color.BgHiBlack, color.Bold).PrintfFunc()
	Ack   = color.New(color.FgHiBlack, color.BgHiGreen, color.Bold).PrintfFunc()
	Warn  = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
)

func Acknowledge(message, content string) error {
	Info1(content)
	fmt.Print("\n")
	Ack(message + " (press enter)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

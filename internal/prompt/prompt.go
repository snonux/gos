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
	INfo3 = color.New(color.FgHiBlack, color.BgHiGreen, color.Bold).PrintfFunc()
	Ack   = color.New(color.FgBlack, color.BgHiYellow, color.Bold).PrintfFunc()
)

func Acknowledge(messages ...string) error {
	if len(messages) > 1 {
		for _, content := range messages[1:] {
			Info1(content)
			fmt.Print("\n")
		}
	}
	Ack(messages[0] + " (press enter to acknowlege)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

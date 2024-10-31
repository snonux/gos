package prompt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	info1 = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	info2 = color.New(color.FgHiYellow, color.BgHiBlack, color.Bold).PrintfFunc()
	ack   = color.New(color.FgHiBlack, color.BgHiGreen, color.Bold).PrintfFunc()
)

func Acknowledge(message, content string) error {
	info1(content)
	fmt.Print("\n")
	ack(message + " (press enter)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

package prompt

import (
	"bufio"
	"fmt"
	"os"

	"codeberg.org/snonux/gos/internal/colour"
)

func Acknowledge(messages ...string) error {
	if len(messages) > 1 {
		for _, content := range messages[1:] {
			colour.Info2f(content)
			fmt.Print("\n")
		}
	}
	colour.Ackf(messages[0] + " (press enter to acknowlege)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

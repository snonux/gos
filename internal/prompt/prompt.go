package prompt

import (
	"bufio"
	"fmt"
	"os"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/table"
)

func Acknowledge(messages ...string) error {
	if len(messages) > 1 {
		for _, content := range messages[1:] {
			table.New().
				WithBaseColor(colour.AttentionCol).
				WithHeaderColor(colour.AckCol).
				Header(messages[0]).
				TextBox(content).
				MustRender()
		}
	}
	fmt.Printf("  ")
	colour.Ackf("(press enter to acknowlege)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

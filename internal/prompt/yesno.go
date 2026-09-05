package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/snonux/gos/internal/colour"
)

// YesNo asks a yes/no question on stdin. It returns true for y/yes and false for n/no.
func YesNo(question string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("  ")
		colour.Ackf("%s (y=yes/n=no):", question)
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("error reading input: %w", err)
		}
		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			fmt.Println("Please respond with one of [yn].")
		}
	}
}

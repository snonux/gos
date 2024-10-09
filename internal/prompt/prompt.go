package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrAborted = errors.New("aborted")

func Yes(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(question, " (y/n): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		switch strings.ToLower(input) {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}

package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var ErrAborted = errors.New("aborted")
var contentSprintf = color.New(color.FgCyan, color.BgBlue, color.Bold).SprintFunc()
var dangerSprintf = color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()

func YesWithContent(question, content string) bool {
	fmt.Print(contentSprintf(content))
	fmt.Print("\n")
	return Yes(question)
}

func Yes(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s ", dangerSprintf(fmt.Sprintf("%s (y/n):", question)))
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

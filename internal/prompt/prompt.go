package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	ErrAborted   = errors.New("aborted")
	contentColor = color.New(color.FgCyan, color.BgBlue, color.Bold).SprintFunc()
	dangerColor  = color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
)

type PromptSelection int

const (
	Unknown PromptSelection = iota
	Yes
	No
	Edit
)

func DoYouWantThis(question, content string) PromptSelection {
	fmt.Print(contentColor(content))
	fmt.Print("\n")
	return whatNow(question)
}

func whatNow(question string) PromptSelection {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s ", dangerColor(fmt.Sprintf("%s (y=yes/n=no/e=edit):", question)))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		switch strings.ToLower(input) {
		case "y", "yes":
			return Yes
		case "n", "no":
			return No
		case "e", "edit":
			return Edit
		default:
			fmt.Println("Please enter 'y' or 'n' or 'e'.")
		}
	}
}

package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var (
	ErrAborted = errors.New("aborted")
	// TODO: Add edit functionality. 1. configure EDITOR, 2. fork EDITOR process on the given file.
	ErrEditContent = errors.New("edit content")
	contentColor   = color.New(color.FgCyan, color.BgBlue, color.Bold).SprintFunc()
	dangerColor    = color.New(color.FgWhite, color.BgRed, color.Bold).SprintFunc()
)

func DoYouWantThis(question, content string) error {
	fmt.Print(contentColor(content))
	fmt.Print("\n")
	return whatNow(question)
}

func whatNow(question string) error {
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
			return nil
		case "n", "no":
			return ErrAborted
		case "e", "edit":
			return ErrEditContent
		default:
			fmt.Println("Please enter 'y' or 'n' or 'e'.")
		}
	}
}

func EditFile(filePath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("EDITOR environment variable is not set")
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

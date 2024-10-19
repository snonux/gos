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
	ErrAborted     = errors.New("aborted")
	ErrEditContent = errors.New("edit content")
	blue           = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	red            = color.New(color.FgWhite, color.BgRed, color.Bold).PrintfFunc()
)

func DoYouWantThis(question, content string) error {
	blue(content)
	fmt.Print("\n")
	return whatNow(question)
}

func Acknowledge(message, content string) error {
	blue(content)
	fmt.Print("\n")
	red(message + " (press enter)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

func whatNow(question string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		red("%s (y=yes/n=no/e=edit):", question)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		switch strings.ToLower(strings.TrimSpace(input)) {
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
	editor, ok := os.LookupEnv("EDITOR")
	if !ok {
		return errors.New("EDITOR environment variable is not set")
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

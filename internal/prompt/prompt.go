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
	ErrDeleteFile  = errors.New("delete file")
	info           = color.New(color.FgCyan, color.BgBlue, color.Bold).PrintfFunc()
	ack            = color.New(color.FgHiBlack, color.BgHiGreen, color.Bold).PrintfFunc()
)

// TODO: Refactor this prompt, including all operations done on the file, like abort, edit, remove, etc.
// TODO: And also don't use error for control flow.
func DoYouWantThis(question, content string) error {
	info(content)
	fmt.Print("\n")
	return whatNow(question)
}

func Acknowledge(message, content string) error {
	info(content)
	fmt.Print("\n")
	ack(message + " (press enter)")
	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		return err
	}
	return nil
}

func whatNow(question string) error {
	reader := bufio.NewReader(os.Stdin)

	for {
		ack("%s (y=yes/n=no/e=edit/d=delete):", question)
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
		case "d", "delete":
			// TODO: Implement
			return ErrDeleteFile
		default:
			fmt.Println("Please enter 'y' or 'n' or 'e' or 'd'.")
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

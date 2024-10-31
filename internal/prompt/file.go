package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"codeberg.org/snonux/gos/internal/oi"
)

var (
	ErrAborted = errors.New("aborted")
	ErrDeleted = errors.New("deleted")
)

func FileAction(question, content, filePath string) error {
	info2(filePath + ":")
	fmt.Print("\n")
	info1(content)
	fmt.Print("\n")
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
			return fmt.Errorf("%w %s", ErrAborted, filePath)
		case "e", "edit":
			if err := EditFile(filePath); err != nil {
				return err
			}
			if content, err = oi.SlurpAndTrim(filePath); err != nil {
				return err
			}
			return FileAction(question, content, filePath)
		case "d", "delete":
			if err := os.Remove(filePath); err != nil {
				return err
			}
			return fmt.Errorf("%w %s", ErrDeleted, filePath)
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

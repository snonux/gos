package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"codeberg.org/snonux/gos/internal/colour"
	"codeberg.org/snonux/gos/internal/oi"
)

var (
	ErrAborted     = errors.New("aborted")
	ErrDeleted     = errors.New("deleted")
	ErrRamdomOther = errors.New("randomOther")
	RandomOption   = true
)

func FileAction(question, content, filePath string, includeRandomOption ...bool) (string, error) {
	colour.Info2fln("%s:", filePath)
	colour.Info2fln("%s", content)
	reader := bufio.NewReader(os.Stdin)

	includeRandom := len(includeRandomOption) > 0 && includeRandomOption[0] == RandomOption
	var randomOption string
	if includeRandom {
		randomOption = "/r=random other"
	}

	for {
		colour.Ackf("%s (y=yes/n=no/e=edit/d=delete%s):", question, randomOption)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %w", err)
		}

		switch strings.ToLower(strings.TrimSpace(input)) {
		case "y", "yes":
			return content, nil
		case "n", "no":
			return content, fmt.Errorf("%w %s", ErrAborted, filePath)
		case "e", "edit":
			if err := EditFile(filePath); err != nil {
				return content, err
			}
			if content, err = oi.SlurpAndTrim(filePath); err != nil {
				return content, err
			}
			return FileAction(question, content, filePath, includeRandomOption...)
		case "d", "delete":
			if err := os.Remove(filePath); err != nil {
				return content, err
			}
			return content, fmt.Errorf("%w %s", ErrDeleted, filePath)
		case "r", "random", "random other":
			if includeRandom {
				return content, fmt.Errorf("%w %s", ErrRamdomOther, filePath)
			}
			fallthrough
		default:
			var r string
			if includeRandom {
				r = "r"
			}
			fmt.Printf("Please respond with one of [yned%s].\n", r)
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

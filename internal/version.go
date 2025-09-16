package internal

import (
	"fmt"

	"codeberg.org/snonux/gos/internal/table"
)

const versionStr = "v1.1.0"

func printVersion() {
	table.New().
		Header("Gos version", "Author", "URL").
		Row(versionStr, "Paul Buetow", "https://codeberg.org/snonux/gos").
		MustRender()
}

func printLogo() {
	raw := `   █████████                    
  ███░░░░░███                  
 ██░░░    ░░░  ██████   █████
░███          ███░░███ ███░░  
░███    █████░███ ░███░░█████ 
░░███  ░░███ ░███ ░███ ░░░░███
 ░░█████████ ░░██████  ██████ 
  ░░░░░░░░░   ░░░░░░  ░░░░░░`

	fmt.Println(raw)
}

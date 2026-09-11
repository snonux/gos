package internal

import (
	"fmt"

	"github.com/snonux/gos/internal/table"
)

const versionStr = "v1.3.4"

func printVersion() {
	table.New().
		Header("Gos version", "Author", "URL").
		Row(versionStr, "Paul Buetow", "https://github.com/snonux/gos").
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

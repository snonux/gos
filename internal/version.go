package internal

import (
	"codeberg.org/snonux/gos/internal/table"
)

const versionStr = "v1.0.0"

func printVersion() {
	table.New().
		Header("Gos version", "Author", "URL").
		Row(versionStr, "Paul Buetow", "https://codeberg.org/snonux/gos").
		MustRender()

	// TODO: Make this work (based on git tag?) From Go 1.24!
	// info, _ := debug.ReadBuildInfo()
	// fmt.Println("Go version:", info.GoVersion)
	// fmt.Println("App version:", info.Main.Version)
}

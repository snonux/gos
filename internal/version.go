package internal

import "codeberg.org/snonux/gos/internal/table"

const versionStr = "v0.0.4"

func printVersion() {
	table.New().
		Header("Gos version", "Author", "URL").
		Row(versionStr, "Paul Buetow", "https://codeberg.org/snonux/gos").
		MustRender()
}

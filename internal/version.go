package internal

import "fmt"

const versionStr = "v0.0.3"

func printVersion() {
	fmt.Printf("This is Gos version %s; (C) by Paul Buetow\n", versionStr)
	fmt.Println("https://codeberg.org/snonux/gos")
}

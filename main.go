package main

import (
	"embed"

	"github.com/LAShZ/go-scaffold/cmd"
)

//go:embed template
var template embed.FS

func main() {
	cmd.Execute(template)
}

package main

import (
	"hjbdev/pvm/commands"
	"hjbdev/pvm/theme"
	"os"
	"runtime"
)

func main() {
	args := os.Args[1:]

	os := runtime.GOOS
	arch := runtime.GOARCH

	if os != "windows" {
		theme.Error("pvm currently only works on Windows.")
		return
	}

	if arch != "amd64" {
		theme.Error("pvm currently only works on 64-bit systems.")
		return
	}

	if len(args) == 0 {
		commands.Help(false)
		return
	}

	switch args[0] {
	case "h", "help":
		commands.Help(false)
	case "l", "ls", "list":
		commands.List()
	case "lr", "ls remote", "ls-remote", "list remote", "list-remote":
		commands.ListRemote()
	case "p", "path":
		commands.Path()
	case "i", "install":
		commands.Install(args)
	case "u", "use":
		commands.Use(args[1:])
	default:
		commands.Help(true)
	}
}

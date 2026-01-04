package main

import (
	"todo_cli/cmd"
)

var (
	version = "dev"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}

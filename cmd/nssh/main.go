package main

import (
	"github.com/0x6b/nssh/cmd"
	"os"
)

func main() {
	os.Exit(run())
}

func run() int {
	if err := cmd.RootCmd.Execute(); err != nil {
		return -1
	}
	return 0
}

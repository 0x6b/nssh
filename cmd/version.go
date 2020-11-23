package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Commit = "tip"
var Tag = "development"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s (%s)\n", Tag, Commit)
		},
	}
}

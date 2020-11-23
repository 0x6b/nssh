package cmd

import (
	"fmt"
	"github.com/0x6b/nssh"
	"github.com/spf13/cobra"
	"os"
)

var (
	coverageType string
	profileName  string
	identity     string
	port         int
	duration     int
	client       *nssh.SoracomClient
)

var RootCmd = &cobra.Command{
	Use:   "nssh name",
	Short: "nssh -- SSH client for SORACOM Napter",
}

func init() {
	RootCmd.PersistentFlags().StringVar(&coverageType, "coverage-type", "", "Specify coverage type, \"g\" for Global, \"jp\" for Japan")
	RootCmd.PersistentFlags().StringVar(&profileName, "profile-name", "ssh", "Specify SORACOM CLI profile name")

	var err error
	client, err = nssh.NewSoracomClient(coverageType, profileName)
	if err != nil {
		fmt.Println("failed to create a client: ", err)
		os.Exit(1)
	}

	RootCmd.AddCommand(listCmd())
	RootCmd.AddCommand(connectCmd())
	RootCmd.AddCommand(versionCmd())
}

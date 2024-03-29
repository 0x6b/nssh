package cmd

import (
	"fmt"
	"github.com/0x6b/nssh/models"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func connectCmd() *cobra.Command {
	connectCmd := &cobra.Command{
		Use:     "connect [<user>@]<subscriber name>",
		Aliases: []string{"c"},
		Short:   "Connect to specified subscriber via SSH.",
		Long:    "Create port mappings for specified subscriber and connect via SSH. If <user>@ is not specified, \"pi\" will be used as default. Quote with \" if name contains spaces or special characters.",
		Args:    cobra.RangeArgs(1, 1),
		Run: func(cmd *cobra.Command, args []string) {
			login, name := parseArg(args[0])

			fmt.Printf("nssh: search subscribers named \"%s\"\n", name)
			onlineSIMs, err := client.FindOnlineSIMsByName(name)
			if err != nil || len(onlineSIMs) == 0 {
				fmt.Printf("nssh: → failed to find online subscribers named \"%s\"\n", name)
				os.Exit(1)
			}

			if len(onlineSIMs) > 1 {
				fmt.Printf("nssh: → cannot create port mapping as there are multiple subscribers named \"%s\"\n", name)
				for _, s := range onlineSIMs {
					fmt.Printf("nssh: - %s\n", s)
				}
				os.Exit(1)
			}

			sim := onlineSIMs[0]
			fmt.Printf("nssh: → found SIM %s\n", sim)

			fmt.Printf("nssh: search existing port mappings for %s:%d\n", sim.ID, port)
			var portMapping *models.PortMapping

			available, err := client.FindAvailablePortMappingsForSIM(sim, port)
			if err != nil || len(available) == 0 {
				fmt.Printf("nssh: → no existing port mapping for %s:%d, creating\n", sim.ID, port)
				portMapping, err = client.CreatePortMappingForSIM(sim, port, duration)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				portMapping = &available[0]
				fmt.Printf("nssh: → found available port mapping:\n%s\n", portMapping)
			}

			fmt.Printf("nssh: connect to %s@%s:%d using the port mapping\n", login, sim.ID, port)
			fmt.Println(strings.Repeat("-", 40))
			err = client.Connect(login, identity, portMapping)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	connectCmd.Flags().StringVarP(&identity, "identity", "i", "", "Specify a path to file from which the identity for public key authentication is read")
	connectCmd.Flags().IntVarP(&port, "port", "p", 22, "Specify port number to connect")
	connectCmd.Flags().IntVarP(&duration, "duration", "d", 60, "Specify session duration in minutes")
	return connectCmd
}

func parseArg(arg string) (string, string) {
	login := "pi"
	var name string

	if strings.Contains(arg, "@") {
		s := strings.SplitN(arg, "@", 2)
		if s[0] != "" {
			login = s[0]
			name = s[1]
		} else {
			name = s[1]
		}
	} else {
		name = arg
	}
	return login, name
}

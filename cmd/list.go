package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func listCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:     "list [subscriber name]",
		Aliases: []string{"l"},
		Short:   "List port mappings for specified subscriber. If no subscriber name is specified, list all port mappings.",
		Args:    cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				portMappings, err := client.ListPortMappings()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				for _, pm := range portMappings {
					subscriber, err := client.GetSubscriber(pm.Destination.Imsi)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println(subscriber)
					fmt.Println(pm)
				}
				return
			}

			subscribers, err := client.FindSubscribersByName(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, s := range subscribers {
				portMappings, err := client.FindPortMappingsForSubscriber(s)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if len(portMappings) > 0 {
					fmt.Println(s)
					for i, pm := range portMappings {
						fmt.Printf("#%d:\n", i+1)
						fmt.Println(pm)
					}
				} else {
					fmt.Printf("no port mapping for %s\n", s)
				}
			}
		},
	}

	return listCmd
}

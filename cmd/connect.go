package cmd

import (
	"fmt"
	"github.com/0x6b/nssh"
	"github.com/spf13/cobra"
	"net"
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
			onlineSubscribers, err := findOnlineSubscribers(name)
			if err != nil || len(onlineSubscribers) == 0 {
				fmt.Printf("nssh: → failed to find online subscribers named \"%s\"\n", name)
				os.Exit(1)
			}

			if len(onlineSubscribers) > 1 {
				fmt.Printf("nssh: → cannot create port mapping as there are multiple subscribers named \"%s\"\n", name)
				for _, s := range onlineSubscribers {
					fmt.Printf("nssh: - %s\n", s)
				}
				os.Exit(1)
			}

			subscriber := onlineSubscribers[0]
			fmt.Printf("nssh: → found subscriber %s\n", subscriber)

			fmt.Printf("nssh: search existing port mappings for %s:%d\n", subscriber.Imsi, port)
			var portMapping *nssh.PortMapping

			available, err := findPortMappings(subscriber, port)
			if err != nil || len(available) == 0 {
				fmt.Printf("nssh: → no existing port mapping for %s:%d, creating\n", subscriber.Imsi, port)
				portMapping, err = client.CreatePortMappingsForSubscriber(subscriber, port, duration)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Println("nssh: → found available port mapping")
				portMapping = &available[0]
			}

			fmt.Printf("nssh: connect to %s:%d using following port mapping:\n%s\n", subscriber.Imsi, port, portMapping)
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

func findOnlineSubscribers(name string) ([]nssh.Subscriber, error) {
	subscribers, err := client.FindSubscribersByName(name)
	if err != nil {
		return nil, err
	}

	var onlineSubscribers []nssh.Subscriber
	for _, s := range subscribers {
		if s.SessionStatus.Online {
			onlineSubscribers = append(onlineSubscribers, s)
		}
	}
	return onlineSubscribers, nil
}

func findPortMappings(subscriber nssh.Subscriber, port int) ([]nssh.PortMapping, error) {
	portMappings, err := client.FindPortMappingsForSubscriber(subscriber)
	if err != nil {
		return nil, err
	}

	var currentPortMappings []nssh.PortMapping
	var availablePortMappings []nssh.PortMapping

	for _, pm := range portMappings {
		if pm.Destination.Port == port {
			currentPortMappings = append(currentPortMappings, pm)
		}
	}

	if len(currentPortMappings) > 0 {
		fmt.Printf("nssh: → found %d port mapping(s) for %s:%d\n", len(currentPortMappings), subscriber.Imsi, port)
		ip, err := nssh.GetIP()

		// search port mappings which allows being connected from current IP address
		if err == nil { // ignore ifconfig.co error
			fmt.Printf("nssh: → current IP address is %s\n", ip)
			for _, pm := range currentPortMappings {
				for _, r := range pm.Source.IPRanges {
					_, ipNet, err := net.ParseCIDR(r)
					if err == nil {
						if ipNet.Contains(ip) {
							availablePortMappings = append(availablePortMappings, pm)
						}
					}
				}
			}
		}
	}
	return availablePortMappings, nil
}

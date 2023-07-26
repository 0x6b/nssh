package cmd

import (
	"fmt"
	"github.com/0x6b/nssh"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list   list.Model
	choice *nssh.Subscriber
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch pressed := msg.String(); pressed {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			s, ok := m.list.SelectedItem().(nssh.Subscriber)
			if ok {
				m.choice = &s
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func (m model) Choice() *nssh.Subscriber {
	return m.choice
}

var login string

func interactiveCmd() *cobra.Command {
	interactiveCmd := &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Short:   "Interactive mode",
		Run: func(cmd *cobra.Command, args []string) {
			subscribers, err := client.FindOnlineSubscribers()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			items := make([]list.Item, 0)

			for _, s := range subscribers {
				if s.Imsi != "" && s.Subscription != "" && s.Type != "" {
					items = append(items, s)
				}
			}

			delegate := list.NewDefaultDelegate()
			delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#34cdd7")).Faint(true)
			delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#34cdd7"))
			delegate.Styles.FilterMatch.Foreground(lipgloss.Color("#34cdd7"))

			m := model{
				list: list.New(items, delegate, 0, 0),
			}
			m.list.Title = "Online Subscribers"
			m.list.Styles.Title = lipgloss.NewStyle().Background(lipgloss.Color("#34cdd7")).Foreground(lipgloss.Color("0")).Bold(true)

			p := tea.NewProgram(m, tea.WithAltScreen())

			result, err := p.Run()
			if err != nil {
				fmt.Println("could not start program:", err)
				os.Exit(1)
			}

			if subscriber := result.(model).Choice(); subscriber != nil {
				fmt.Printf("nssh: search existing port mappings for %s:%d\n", subscriber.Imsi, port)
				var portMapping *nssh.PortMapping

				available, err := findPortMappings(*subscriber, port)
				if err != nil || len(available) == 0 {
					fmt.Printf("nssh: → no existing port mapping for %s:%d, creating\n", subscriber.Imsi, port)
					portMapping, err = client.CreatePortMappingsForSubscriber(*subscriber, port, duration)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				} else {
					portMapping = &available[0]
					fmt.Printf("nssh: → found available port mapping:\n%s\n", portMapping)
				}

				fmt.Printf("nssh: connect to %s:%d using the port mapping\n", subscriber.Imsi, port)
				fmt.Println(strings.Repeat("-", 40))
				err = client.Connect(login, identity, portMapping)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}

	interactiveCmd.Flags().StringVarP(&login, "login", "u", "pi", "Specify login user name")
	interactiveCmd.Flags().StringVarP(&identity, "identity", "i", "", "Specify a path to file from which the identity for public key authentication is read")
	interactiveCmd.Flags().IntVarP(&port, "port", "p", 22, "Specify port number to connect")
	interactiveCmd.Flags().IntVarP(&duration, "duration", "d", 60, "Specify session duration in minutes")
	return interactiveCmd
}

package main

import (
	"net"

	"github.com/etcdhosts/client-go"

	"github.com/muesli/reflow/indent"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/urfave/cli/v2"
)

func addCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a DNS record",
		UsageText: "add HOST IP",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 || net.ParseIP(c.Args().Get(1)) == nil {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}
			return tea.NewProgram(addMode{ctx: c, stages: 5}).Start()
		},
	}
}

type addMode struct {
	ctx      *cli.Context
	hc       *client.Client
	hf       *client.HostsFile
	stages   int
	stageMsg string
	progress float64
	loaded   bool
}

func (m addMode) Init() tea.Cmd {
	return nil
}

// The main view, which just calls the approprate sub-view
func (m addMode) View() string {
	return indent.String("\n"+m.stageMsg, 2) + indent.String("\n"+progressbar(cmdProgressBarWidth, m.progress)+"%"+"\n\n", 2)
}

// Main update function.
func (m addMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	var err error
	var stage int
	switch msg.(type) {
	case int:
		stage = msg.(int)
	default:
		return m, func() tea.Msg { return stage }
	}

	switch stage {
	case 0:
		m.stageMsg = "Create etcdhosts client..."
		m.hc, err = createClient(m.ctx)
		if err != nil {
			m.stageMsg = "ERROR: " + err.Error()
			return m, tea.Quit
		}
	case 1:
		m.stageMsg = "Read hostsfile from etcd..."
		m.hf, err = m.hc.ReadHostsFile()
		if err != nil {
			m.stageMsg = "ERROR: " + err.Error()
			return m, tea.Quit
		}
	case 2:
		m.stageMsg = "Add host to hostsfile..."
		err = m.hf.AddHost(m.ctx.Args().Get(0), m.ctx.Args().Get(1))
		if err != nil {
			m.stageMsg = "ERROR: " + err.Error()
			return m, tea.Quit
		}
	case 3:
		m.stageMsg = "Put hostsfile into etcd..."
		err = m.hc.PutHostsFile(m.hf)
		if err != nil {
			m.stageMsg = "ERROR: " + err.Error()
			return m, tea.Quit
		}
	case 4:
		m.stageMsg = "DNS records add success..."
	}

	stage++
	if !m.loaded {
		m.progress += float64(1) / float64(m.stages)
		if m.progress >= 1 {
			m.progress = 1
			m.loaded = true
			return m, tea.Quit
		}
	}

	return m, func() tea.Msg { return stage }
}

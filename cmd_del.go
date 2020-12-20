package main

import (
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etcdhosts/client-go"
	"github.com/muesli/reflow/indent"

	"github.com/urfave/cli/v2"
)

func delCmd() *cli.Command {
	return &cli.Command{
		Name:      "del",
		Usage:     "Delete a DNS record",
		UsageText: "del HOST [IP]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "purge",
				Value: false,
				Usage: "delete all records of a given host",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("purge") && c.NArg() != 1 {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}

			if !c.Bool("purge") && (c.NArg() != 2 || net.ParseIP(c.Args().Get(1)) == nil) {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}

			return tea.NewProgram(delMode{ctx: c, stages: 5}).Start()
		},
	}
}

type delMode struct {
	ctx      *cli.Context
	hc       *client.Client
	hf       *client.HostsFile
	stages   int
	stageMsg string
	progress float64
	loaded   bool
}

func (m delMode) Init() tea.Cmd {
	return nil
}

// The main view, which just calls the approprate sub-view
func (m delMode) View() string {
	return indent.String("\n"+m.stageMsg, 2) + indent.String("\n"+progressbar(cmdProgressBarWidth, m.progress)+"%"+"\n\n", 2)
}

// Main update function.
func (m delMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.stageMsg = makeInfo("Create etcdhosts client...")
		m.hc, err = createClient(m.ctx)
		if err != nil {
			m.stageMsg = makeError("ERROR: " + err.Error())
			return m, tea.Quit
		}
	case 1:
		m.stageMsg = makeInfo("Read hostsfile from etcd...")
		m.hf, err = m.hc.ReadHostsFile()
		if err != nil {
			m.stageMsg = makeError("ERROR: " + err.Error())
			return m, tea.Quit
		}
	case 2:
		m.stageMsg = makeInfo("Delete host from hostsfile...")
		if m.ctx.Bool("purge") {
			m.hf.PurgeHost(m.ctx.Args().Get(0))
		} else {
			err = m.hf.DelHost(m.ctx.Args().Get(0), m.ctx.Args().Get(1))
			if err != nil {
				m.stageMsg = makeError("ERROR: " + err.Error())
				return m, tea.Quit
			}
		}
	case 3:
		m.stageMsg = makeInfo("Put hostsfile into etcd...")
		err = m.hc.PutHostsFile(m.hf)
		if err != nil {
			m.stageMsg = makeError("ERROR: " + err.Error())
			return m, tea.Quit
		}
	case 4:
		m.stageMsg = makeInfo("DNS records delete success...")
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

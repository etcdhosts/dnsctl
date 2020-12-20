package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/etcdhosts/client-go"
	"github.com/muesli/reflow/indent"

	"github.com/urfave/cli/v2"
)

func restoreCmd() *cli.Command {
	return &cli.Command{
		Name:      "restore",
		Usage:     "Restore dns records from hosts file",
		UsageText: "restore FILE_PATH",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "backup",
				Aliases: []string{"b"},
				Value:   true,
				Usage:   "back up the original hosts file when restoring",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}
			f, err := os.Open(c.Args().Get(0))
			if err != nil {
				return err
			}
			_ = f.Close()

			return tea.NewProgram(restoreMode{ctx: c, stages: 4}).Start()
		},
	}
}

type restoreMode struct {
	ctx      *cli.Context
	hc       *client.Client
	stages   int
	stageMsg string
	progress float64
	loaded   bool
}

func (m restoreMode) Init() tea.Cmd {
	return nil
}

// The main view, which just calls the approprate sub-view
func (m restoreMode) View() string {
	return indent.String("\n"+m.stageMsg, 2) + indent.String("\n"+progressbar(cmdProgressBarWidth, m.progress)+"%"+"\n\n", 2)
}

// Main update function.
func (m restoreMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.stageMsg = makeInfo("Backup hostsfile from etcd...")
		if m.ctx.Bool("backup") {
			hf, err := m.hc.ReadHostsFile()
			if err != nil {
				m.stageMsg = makeError("ERROR: " + err.Error())
				return m, tea.Quit
			}
			err = ioutil.WriteFile(fmt.Sprintf("dnsctl.%d.bak", time.Now().Unix()), []byte(hf.String()), 0644)
			if err != nil {
				m.stageMsg = makeError("ERROR: " + err.Error())
				return m, tea.Quit
			}
		}
	case 2:
		m.stageMsg = makeInfo("Restore from file...")
		f, err := os.Open(m.ctx.Args().Get(0))
		if err != nil {
			return m, tea.Quit
		}
		defer func() { _ = f.Close() }()
		err = m.hc.ForcePutHostsFile(f)
		if err != nil {
			return m, tea.Quit
		}
	case 3:
		m.stageMsg = makeInfo("DNS records restore success...")
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

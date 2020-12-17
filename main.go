package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatal(err)
	}
	app := &cli.App{
		Name:    "dnsctl",
		Usage:   "Command line tool for etcdhosts plugin",
		Version: fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors: []*cli.Author{
			{
				Name:  "mritd",
				Email: "mritd@linux.com",
			},
		},
		Copyright: "Copyright (c) 2020 mritd, All rights reserved.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   filepath.Join(home, ".dnsctl.yaml"),
				Usage:   "dnsctl config",
				EnvVars: []string{"DNSCTL_CONFIG"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Value:   false,
				Usage:   "debug mode",
				EnvVars: []string{"DNSCTL_DEBUG"},
			},
		},
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				logger.SetDevelopment()
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			exampleCmd(),
			addCmd(),
			delCmd(),
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}

func exampleCmd() *cli.Command {
	return &cli.Command{
		Name:  "example",
		Usage: "print example config",
		Action: func(c *cli.Context) error {
			fmt.Println(ExampleConfig())
			return nil
		},
	}
}

func addCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "add a DNS record",
		UsageText: "add HOST IP",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 || net.ParseIP(c.Args().Get(1)) == nil {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHosts()
			if err != nil {
				return err
			}

			err = hf.AddHost(c.Args().Get(0), c.Args().Get(1))
			if err != nil {
				return err
			}
			logger.Info("host add success.")
			return nil
		},
	}
}

func delCmd() *cli.Command {
	return &cli.Command{
		Name:      "del",
		Usage:     "delete a DNS record",
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

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHosts()
			if err != nil {
				return err
			}
			if c.Bool("purge") {
				hf.PurgeHost(c.Args().Get(0))
				return nil
			} else {
				err = hf.DelHost(c.Args().Get(0), c.Args().Get(1))
				if err != nil {
					return err
				}
			}

			logger.Info("host delete success.")
			return nil
		},
	}
}

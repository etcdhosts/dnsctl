package main

import (
	"fmt"
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
		Copyright: "Copyright (c) 2022 mritd, All rights reserved.",
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
			purgeCmd(),
			dumpCmd(),
			restoreCmd(),
			historyCmd(),
			editCmd(),
		},
	}
	if err = app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

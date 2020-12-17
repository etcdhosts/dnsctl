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
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				logger.SetDevelopment()
			}

			return nil
		},
		Commands: []*cli.Command{
			exampleCmd(),
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

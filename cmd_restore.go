package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func restoreCmd() *cli.Command {
	return &cli.Command{
		Name:      "restore",
		Usage:     "Restore dns records from hosts file",
		UsageText: "restore FILE_PATH",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}
			f, err := os.Open(c.Args().Get(0))
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			return hc.ForcePutHostsFile(f)
		},
	}
}

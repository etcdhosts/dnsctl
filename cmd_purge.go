package main

import (
	"github.com/urfave/cli/v2"
)

func purgeCmd() *cli.Command {
	return &cli.Command{
		Name:      "purge",
		Usage:     "Delete all DNS records for a given HOST",
		UsageText: "purge HOST",
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}

			hf.PurgeHost(c.Args().Get(0))

			return hc.PutHostsFile(hf)
		},
	}
}

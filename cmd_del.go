package main

import (
	"net"

	"github.com/urfave/cli/v2"
)

func delCmd() *cli.Command {
	return &cli.Command{
		Name:      "del",
		Usage:     "Delete a DNS record",
		UsageText: "del IP HOST",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 || net.ParseIP(c.Args().Get(0)) == nil {
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

			err = hf.DelHost(c.Args().Get(1), c.Args().Get(0))
			if err != nil {
				return err
			}

			return hc.PutHostsFile(hf)
		},
	}
}

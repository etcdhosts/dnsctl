package main

import (
	"github.com/urfave/cli/v2"
	"net"
)

func addCmd() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "Add a DNS record",
		UsageText: "add IP HOST",
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				cli.ShowCommandHelpAndExit(c, c.Command.Name, 1)
			}

			if net.ParseIP(c.Args().Get(0)) == nil {
				return cli.Exit("Args[0] must be a IP format",1)
			}

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}

			err = hf.AddHost(c.Args().Get(1), c.Args().Get(0))
			if err != nil {
				return err
			}

			return hc.PutHostsFile(hf)
		},
	}
}

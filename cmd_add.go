package main

import (
	"net"

	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

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

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}

			err = hf.AddHost(c.Args().Get(0), c.Args().Get(1))
			if err != nil {
				return err
			}

			err = hc.PutHostsFile(hf)
			if err != nil {
				return err
			}

			logger.Info("DNS record added successfully.")
			return nil
		},
	}
}

package main

import (
	"net"

	"github.com/mritd/logger"
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

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}
			if c.Bool("purge") {
				hf.PurgeHost(c.Args().Get(0))
			} else {
				err = hf.DelHost(c.Args().Get(0), c.Args().Get(1))
				if err != nil {
					return err
				}
			}

			err = hc.PutHostsFile(hf)
			if err != nil {
				return err
			}

			logger.Info("DNS record deleted successfully.")
			return nil
		},
	}
}

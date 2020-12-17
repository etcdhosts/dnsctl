package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/mritd/logger"
	"github.com/urfave/cli/v2"
)

func restoreCmd() *cli.Command {
	return &cli.Command{
		Name:      "restore",
		Usage:     "restore dns records from hosts file",
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
			defer func() { _ = f.Close() }()

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			if c.Bool("backup") {
				hf, err := hc.ReadHostsFile()
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(fmt.Sprintf("dnsctl.%d.bak", time.Now().Unix()), []byte(hf.String()), 0644)
				if err != nil {
					return err
				}
			}

			err = hc.ForcePutHostsFile(f)
			if err != nil {
				return err
			}

			logger.Info("DNS record restored successfully.")
			return nil
		},
	}
}

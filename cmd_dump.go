package main

import (
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli/v2"
)

func dumpCmd() *cli.Command {
	return &cli.Command{
		Name:  "dump",
		Usage: "dump all DNS records",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "dump file storage location",
			},
		},
		Action: func(c *cli.Context) error {
			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}
			if c.String("out") == "" {
				fmt.Println(hf.String())
				return nil
			} else {
				return ioutil.WriteFile(c.String("out"), []byte(hf.String()), 0644)
			}
		},
	}
}

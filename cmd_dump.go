package main

import (
	"fmt"
	"io/ioutil"

	"github.com/etcdhosts/client-go"

	"github.com/urfave/cli/v2"
)

func dumpCmd() *cli.Command {
	return &cli.Command{
		Name:  "dump",
		Usage: "Dump all DNS records",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "dump file storage location",
			},
			&cli.Int64Flag{
				Name:    "revision",
				Aliases: []string{"r"},
				Value:   0,
				Usage:   "dump hosts file from etcd revision",
			},
		},
		Action: func(c *cli.Context) error {
			hc, err := createClient(c)
			if err != nil {
				return err
			}

			var hf *client.HostsFile
			if c.Int64("revision") > 0 {
				hf, err = hc.ReadHostsFileByRevision(c.Int64("revision"))
				if err != nil {
					return err
				}
			} else {
				hf, err = hc.ReadHostsFile()
				if err != nil {
					return err
				}
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

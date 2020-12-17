package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"time"

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
		Before: func(c *cli.Context) error {
			if c.Bool("debug") {
				logger.SetDevelopment()
			}
			return nil
		},
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			exampleCmd(),
			addCmd(),
			delCmd(),
			dumpCmd(),
			restoreCmd(),
		},
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Error(err)
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

func delCmd() *cli.Command {
	return &cli.Command{
		Name:      "del",
		Usage:     "delete a DNS record",
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

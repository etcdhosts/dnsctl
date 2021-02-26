package main

import (
	"fmt"

	"github.com/etcdhosts/client-go"

	"github.com/urfave/cli/v2"
)

func historyCmd() *cli.Command {
	return &cli.Command{
		Name:  "history",
		Usage: "Print hosts change history",
		Action: func(c *cli.Context) error {
			return fmt.Errorf("the history command has not yet been implemented")
		},
	}
}

type historyMode struct {
	ctx *cli.Context
	hc  *client.Client
	hf  *client.HostsFile
}

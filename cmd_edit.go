package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func editCmd() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "Edit the hosts file with an editor",
		Action: func(c *cli.Context) error {
			return fmt.Errorf("the edit command has not yet been implemented")
		},
	}
}

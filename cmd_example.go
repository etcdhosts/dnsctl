package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

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

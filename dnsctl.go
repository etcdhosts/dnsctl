package main

import (
	"github.com/etcdhosts/client-go"
	"github.com/urfave/cli/v2"
)

func createClient(c *cli.Context) (*client.Client, error) {
	var cfg Config
	err := cfg.LoadFrom(c.String("config"))
	if err != nil {
		return nil, err
	}

	return client.NewClient(cfg.Config)
}

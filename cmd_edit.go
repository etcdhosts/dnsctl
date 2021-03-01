package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/mritd/logger"

	"github.com/etcdhosts/client-go"

	"github.com/urfave/cli/v2"
)

func editCmd() *cli.Command {
	return &cli.Command{
		Name:  "edit",
		Usage: "Edit the hosts file with an editor",
		Action: func(c *cli.Context) error {
			f, err := ioutil.TempFile("", "dnsctl")
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
				_ = os.Remove(f.Name())
			}()

			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hf, err := hc.ReadHostsFile()
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(f, hf.String())
			if err != nil {
				return err
			}

			// get os editor
			editor := "vim"
			if runtime.GOOS == "windows" {
				editor = "notepad"
			}
			if v := os.Getenv("VISUAL"); v != "" {
				editor = v
			} else if e := os.Getenv("EDITOR"); e != "" {
				editor = e
			}

			cmd := exec.Command(editor, f.Name())
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return err
			}

			raw, err := ioutil.ReadFile(f.Name())
			if err != nil {
				return err
			}
			hm := client.Parse2Map(bytes.NewReader(raw))
			if hm.String() == hf.String() {
				logger.Info("DNS records not change.")
				return nil
			} else {
				logger.Info("DNS records updated.")
				return hc.ForcePutHostsFile(bytes.NewReader(raw))
			}
		},
	}
}

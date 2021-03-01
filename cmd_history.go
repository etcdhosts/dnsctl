package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/urfave/cli/v2"
)

const historyTpl = `+-------------+------------------+
|   VERSION   |   MOD REVISION   |
+-------------+------------------+
{{range $index,$hf := .}}|   {{$hf.Version | printf "%-7d"}}   |   {{ $hf.ModRevision | printf "%-12d   |\n"}}{{end}}+--------------------------------+
`

func historyCmd() *cli.Command {
	return &cli.Command{
		Name:  "history",
		Usage: "Print hosts change history",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Value:   "",
				Usage:   "dump all version to dir",
			},
		},
		Action: func(c *cli.Context) error {
			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hfs, err := hc.GetHostsFileHistory()
			if err != nil {
				return err
			}
			if dir := c.String("out"); dir != "" {
				info, err := os.Stat(dir)
				if err != nil {
					if os.IsNotExist(err) {
						err = os.MkdirAll(dir, 0755)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				} else {
					if !info.IsDir() {
						return fmt.Errorf("[%s] is not a dir", dir)
					}
				}

				for _, hf := range hfs {
					f, err := os.OpenFile(filepath.Join(dir, fmt.Sprintf("dnsctl_%d", hf.Version())), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
					if err != nil {
						return err
					}
					_, err = fmt.Fprintln(f, hf.String())
					if err != nil {
						return err
					}
					_ = f.Close()
				}
				return nil
			} else {
				tpl, err := template.New("history").Parse(historyTpl)
				if err != nil {
					return err
				}
				return tpl.Execute(os.Stdout, hfs)
			}
		},
	}
}

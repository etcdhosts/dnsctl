package main

import (
	"os"
	"text/template"

	"github.com/urfave/cli/v2"
)

const historyTpl = `+-------------+------------------+
|   VERSION   |   MOD REVISION   |
+-------------+------------------+
{{range $index,$hf := .}}|   {{$hf.Version | printf "%-7d"}}   |   {{ $hf.ModRevision | printf "%-12d"}}   |{{end}}
+--------------------------------+
`

func historyCmd() *cli.Command {
	return &cli.Command{
		Name:  "history",
		Usage: "Print hosts change history",
		Action: func(c *cli.Context) error {
			hc, err := createClient(c)
			if err != nil {
				return err
			}

			hfs, err := hc.GetHostsFileHistory()
			if err != nil {
				return err
			}

			tpl, err := template.New("history").Parse(historyTpl)
			if err != nil {
				return err
			}
			return tpl.Execute(os.Stdout, hfs)
		},
	}
}

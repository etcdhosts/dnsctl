package cmd

import (
	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/output"
)

var listOutput string
var listRevision int64

// listCmd represents the list command.
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all DNS records",
	Long: `List all DNS records stored in etcd.

Output formats:
  hosts - hosts file format (default)
  json  - JSON format
  yaml  - YAML format

Example:
  dnsctl list
  dnsctl list -o json
  dnsctl list -o yaml
  dnsctl list -r 12345`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listOutput, "output", "o", "hosts", "output format: hosts, json, yaml")
	listCmd.Flags().Int64VarP(&listRevision, "revision", "r", 0, "read from specific etcd revision")
}

func runList(cmd *cobra.Command, args []string) error {
	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	var hosts any
	if listRevision > 0 {
		hosts, err = cli.ReadRevision(listRevision)
	} else {
		hosts, err = cli.Read()
	}
	if err != nil {
		return err
	}

	return output.Print(hosts, output.Format(listOutput))
}

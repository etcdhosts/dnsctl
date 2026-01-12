package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
	defer cli.Close()

	var hosts any
	if listRevision > 0 {
		hosts, err = cli.ReadRevision(listRevision)
	} else {
		hosts, err = cli.Read()
	}
	if err != nil {
		return err
	}

	return outputHosts(hosts, listOutput)
}

func outputHosts(hosts any, format string) error {
	type stringer interface {
		String() string
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(hosts)
	case "yaml":
		return yaml.NewEncoder(os.Stdout).Encode(hosts)
	default:
		if s, ok := hosts.(stringer); ok {
			fmt.Print(s.String())
		}
		return nil
	}
}

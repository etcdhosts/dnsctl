package cmd

import (
	"fmt"
	"strings"

	client "github.com/etcdhosts/client-go/v2"
	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/output"
)

var historyLimit int

// historyCmd represents the history command.
var historyCmd = &cobra.Command{
	Use:   "history [domain]",
	Short: "List all available versions",
	Long: `List all historical versions of DNS records stored in etcd.

For single-key mode:
  Shows all versions of the hosts data.

For per-host mode:
  Without argument: lists all domains
  With domain argument: shows history for that domain

Use 'dnsctl list -r REVISION' to view a specific version.

Example:
  dnsctl history
  dnsctl history -n 10
  dnsctl history example.com`,
	RunE: runHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)

	historyCmd.Flags().IntVarP(&historyLimit, "limit", "n", 0, "limit number of versions to show (0 = all)")
}

func runHistory(cmd *cobra.Command, args []string) error {
	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	mode, err := cli.Mode()
	if err != nil {
		return err
	}

	if mode == client.ModePerHost {
		if len(args) == 0 {
			return listDomains(cli)
		}
		return showDomainHistory(cli, args[0])
	}

	return showHistory(cli)
}

func listDomains(cli *client.Client) error {
	domains, err := cli.ListDomains()
	if err != nil {
		return err
	}

	if len(domains) == 0 {
		fmt.Println("No domains found.")
		return nil
	}

	fmt.Println("Domains:")
	for _, d := range domains {
		fmt.Printf("  %s\n", strings.TrimSuffix(d, "."))
	}
	fmt.Printf("\nTotal: %d domains\n", len(domains))
	fmt.Println("Use 'dnsctl history <domain>' to view domain history.")
	return nil
}

func showDomainHistory(cli *client.Client, domain string) error {
	history, err := cli.HistoryHost(domain)
	if err != nil {
		return err
	}

	if len(history) == 0 {
		fmt.Printf("No history for domain: %s\n", domain)
		return nil
	}

	fmt.Printf("History for %s:\n\n", domain)
	output.PrintHistoryTable(history, historyLimit)
	return nil
}

func showHistory(cli *client.Client) error {
	history, err := cli.History()
	if err != nil {
		return err
	}

	if len(history) == 0 {
		fmt.Println("No history available.")
		return nil
	}

	output.PrintHistoryTable(history, historyLimit)
	return nil
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// purgeCmd represents the purge command.
var purgeCmd = &cobra.Command{
	Use:   "purge HOSTNAME",
	Short: "Delete all DNS records for a hostname",
	Long: `Delete all DNS records associated with a hostname.

This removes all IP mappings for the specified hostname.

Example:
  dnsctl purge example.com`,
	Args: cobra.ExactArgs(1),
	RunE: runPurge,
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}

func runPurge(cmd *cobra.Command, args []string) error {
	hostname := args[0]

	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	hosts, err := cli.Read()
	if err != nil {
		return err
	}

	hosts.Purge(hostname)

	if err := cli.Write(hosts); err != nil {
		return err
	}

	fmt.Printf("Purged: %s\n", hostname)
	return nil
}

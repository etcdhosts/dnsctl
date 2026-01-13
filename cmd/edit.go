package cmd

import (
	"errors"
	"fmt"

	client "github.com/etcdhosts/client-go/v2"
	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/editor"
	"github.com/etcdhosts/dnsctl/v2/internal/output"
)

// editCmd represents the edit command.
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit DNS records with system editor",
	Long: `Open DNS records in system editor for editing.

The editor is determined by:
  1. $EDITOR environment variable
  2. $VISUAL environment variable
  3. Default: vi

Example:
  dnsctl edit
  EDITOR=nano dnsctl edit`,
	RunE: runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	hosts, err := cli.Read()
	if err != nil {
		return err
	}

	result, err := editor.Edit(hosts.String())
	if err != nil {
		return err
	}

	if !result.Modified {
		fmt.Println("No changes made.")
		return nil
	}

	// Use strict parsing to detect errors
	parseResult := client.ParseRecordsStrict(result.Content)
	if parseResult.HasErrors() {
		fmt.Printf("Error: found %d invalid record(s):\n", len(parseResult.Errors))
		for _, e := range parseResult.Errors {
			fmt.Printf("  - %s\n", e.String())
		}
		fmt.Println("\nNo changes were saved. Please fix the errors and try again.")
		return nil
	}

	newHosts, warnings := dedupeRecords(parseResult.Records)

	if len(warnings) > 0 {
		fmt.Printf("Warning: removed %d duplicate record(s):\n", len(warnings))
		for _, warn := range warnings {
			fmt.Printf("  - %s\n", warn)
		}
	}

	newHosts.SetModified(hosts.Modified())
	if err := cli.ForceWrite([]byte(newHosts.String())); err != nil {
		return err
	}

	fmt.Printf("Updated %d records.\n", newHosts.Len())
	return nil
}

// dedupeRecords creates a new Hosts with duplicates removed.
// Returns the deduplicated Hosts and warning messages for removed duplicates.
func dedupeRecords(records []client.Record) (*client.Hosts, []string) {
	newHosts := client.NewHosts()
	var warnings []string

	for _, r := range records {
		if err := newHosts.Add(r); err != nil {
			if errors.Is(err, client.ErrDuplicateRecord) {
				warnings = append(warnings,
					fmt.Sprintf("%s -> %s %s (duplicate, removed)", r.Hostname, r.IP, output.FormatRecordAttrs(r)))
			}
		}
	}

	return newHosts, warnings
}

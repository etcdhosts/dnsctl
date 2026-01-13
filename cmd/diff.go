package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/diff"
)

// diffCmd represents the diff command.
var diffCmd = &cobra.Command{
	Use:   "diff <revision1> <revision2>",
	Short: "Compare two versions of DNS records",
	Long: `Compare two versions of DNS records and show differences.

Shows added lines in green and removed lines in red.
Use 'dnsctl history' to list available revisions.

Example:
  dnsctl diff 100 200
  dnsctl diff 100 0      # compare revision 100 with current`,
	Args: cobra.ExactArgs(2),
	RunE: runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	var rev1, rev2 int64
	if _, err := fmt.Sscanf(args[0], "%d", &rev1); err != nil {
		return fmt.Errorf("invalid revision: %s", args[0])
	}
	if _, err := fmt.Sscanf(args[1], "%d", &rev2); err != nil {
		return fmt.Errorf("invalid revision: %s", args[1])
	}

	cli, err := newClient()
	if err != nil {
		return err
	}
	defer func() { _ = cli.Close() }()

	// Get first version
	hosts1, err := cli.ReadRevision(rev1)
	if err != nil {
		return fmt.Errorf("failed to read revision %d: %w", rev1, err)
	}

	// Get second version (0 means current)
	hosts2, err := cli.ReadRevision(rev2)
	if err != nil {
		return fmt.Errorf("failed to read revision %d: %w", rev2, err)
	}

	// Get string representations
	str1 := hosts1.String()
	str2 := hosts2.String()

	if str1 == str2 {
		fmt.Println("No differences found.")
		return nil
	}

	// Print header
	fmt.Printf("%s--- revision %d%s\n", diff.ColorYellow, rev1, diff.ColorReset)
	fmt.Printf("%s+++ revision %d%s\n", diff.ColorYellow, rev2, diff.ColorReset)
	fmt.Println()

	// Compute and print diff
	diff.PrintUnified(str1, str2)
	return nil
}

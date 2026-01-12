package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/config"
)

// exampleCmd represents the example command.
var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Print example config",
	Long:  `Print an example configuration file for dnsctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.Example())
	},
}

func init() {
	rootCmd.AddCommand(exampleCmd)
}

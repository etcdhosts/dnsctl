// Package cmd contains all CLI commands.
package cmd

import (
	"os"
	"path/filepath"

	client "github.com/etcdhosts/client-go/v2"
	"github.com/spf13/cobra"

	"github.com/etcdhosts/dnsctl/v2/internal/config"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "dnsctl",
	Short: "Command line tool for etcdhosts DNS management",
	Long: `dnsctl is a CLI tool for managing DNS records stored in etcd.

It provides commands to edit, list, compare, and manage DNS records
that are used by the etcdhosts CoreDNS plugin.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	home, _ := os.UserHomeDir()
	defaultConfig := filepath.Join(home, ".dnsctl.yaml")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultConfig, "config file path")
}

// SetVersion sets the version string for --version flag.
func SetVersion(version string) {
	rootCmd.Version = version
}

// newClient creates a new etcdhosts client from config.
func newClient() (*client.Client, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}
	return client.NewClient(cfg.ToClientConfig())
}

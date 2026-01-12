package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

// delCmd represents the del command.
var delCmd = &cobra.Command{
	Use:   "del IP HOSTNAME",
	Short: "Delete a DNS record",
	Long: `Delete a specific DNS record by IP and hostname.

Example:
  dnsctl del 192.168.1.1 example.com`,
	Args: cobra.ExactArgs(2),
	RunE: runDel,
}

func init() {
	rootCmd.AddCommand(delCmd)
}

func runDel(cmd *cobra.Command, args []string) error {
	ipStr := args[0]
	hostname := args[1]

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	cli, err := newClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	hosts, err := cli.Read()
	if err != nil {
		return err
	}

	if err := hosts.Del(hostname, ip); err != nil {
		return err
	}

	if err := cli.Write(hosts); err != nil {
		return err
	}

	fmt.Printf("Deleted: %s -> %s\n", hostname, ipStr)
	return nil
}

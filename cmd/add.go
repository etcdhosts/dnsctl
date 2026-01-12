package cmd

import (
	"fmt"
	"net"

	client "github.com/etcdhosts/client-go/v2"
	"github.com/spf13/cobra"
)

var addWeight int
var addTTL uint32

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add IP HOSTNAME",
	Short: "Add a DNS record",
	Long: `Add a new DNS record mapping an IP address to a hostname.

Example:
  dnsctl add 192.168.1.1 example.com
  dnsctl add 192.168.1.1 example.com -w 10
  dnsctl add 192.168.1.1 example.com -t 60`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().IntVarP(&addWeight, "weight", "w", 1, "weight for load balancing (1-10000)")
	addCmd.Flags().Uint32VarP(&addTTL, "ttl", "t", 0, "TTL in seconds")
}

func runAdd(cmd *cobra.Command, args []string) error {
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

	record := client.Record{
		Hostname: hostname,
		IP:       ip,
		Weight:   addWeight,
		TTL:      addTTL,
	}

	if err := hosts.Add(record); err != nil {
		return err
	}

	if err := cli.Write(hosts); err != nil {
		return err
	}

	fmt.Printf("Added: %s -> %s\n", hostname, ipStr)
	return nil
}

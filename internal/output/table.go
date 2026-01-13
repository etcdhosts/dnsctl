package output

import (
	"fmt"

	client "github.com/etcdhosts/client-go/v2"
)

// PrintHistoryTable prints a table of hosts history.
func PrintHistoryTable(history []*client.Hosts, limit int) {
	total := len(history)

	// Apply limit if specified
	if limit > 0 && limit < len(history) {
		history = history[:limit]
	}

	fmt.Printf("%-14s  %-10s  %-8s  %s\n", "REVISION", "VERSION", "RECORDS", "MODIFIED")
	fmt.Println("--------------  ----------  --------  --------------------")

	for i, h := range history {
		marker := ""
		if i == 0 {
			marker = " (latest)"
		}

		modified := "-"
		if !h.Modified().IsZero() {
			modified = h.Modified().Local().Format("2006-01-02 15:04:05")
		}

		fmt.Printf("%-14d  %-10d  %-8d  %s%s\n",
			h.ModRevision(),
			h.Version(),
			h.Len(),
			modified,
			marker,
		)
	}

	fmt.Printf("\nShowing %d of %d versions\n", len(history), total)
	fmt.Println("Use 'dnsctl list -r REVISION' to view a specific version.")
}

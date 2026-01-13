package main

import (
	"fmt"

	"github.com/etcdhosts/dnsctl/v2/cmd"
)

var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "unknown"
)

func main() {
	cmd.SetVersion(fmt.Sprintf("%s %s %s", version, buildDate, commit))
	cmd.Execute()
}

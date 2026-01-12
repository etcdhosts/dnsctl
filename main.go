package main

import "github.com/etcdhosts/dnsctl/v2/cmd"

var (
	version   = "dev"
	buildDate = "unknown"
	commit    = "unknown"
)

func main() {
	cmd.SetVersionInfo(cmd.VersionInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	})
	cmd.Execute()
}

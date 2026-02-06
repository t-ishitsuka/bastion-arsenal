package main

import "github.com/arsenal/internal/cli"

// Version information (overwritten at build time with -ldflags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Arsenal のエントリポイント
func main() {
	cli.SetVersion(Version, Commit, BuildDate)
	cli.Execute()
}

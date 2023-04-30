package main

import (
	cli "github.com/jawher/mow.cli"
)

func cmdCompare(cmd *cli.Cmd) {
	cmd.Spec = "SCAN..."
	var (
		scans = cmd.StringsArg("SCAN", nil, "Scans to compare")
	)

	cmd.Action = func() {
		compare(*scans)
	}

}

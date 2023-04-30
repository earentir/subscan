package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

func cmdList(cmd *cli.Cmd) {
	cmd.Spec = "[DIR]"
	var (
		dir = cmd.StringArg("DIR", ".", "Directory to list")
	)

	cmd.Action = func() {
		matchedFiles, err := getMatchingFiles(*dir)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Found", len(matchedFiles), "scans")
		fmt.Println()

		groupedFiles := groupFilesByIPSubnet(matchedFiles)

		printGroupedFiles(groupedFiles, *dir)
	}

}

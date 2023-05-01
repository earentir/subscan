package main

import (
	"fmt"
	"os"

	cli "github.com/jawher/mow.cli"
)

var (
	appversion string = "0.1.5"
)

func main() {
	app := cli.App("subscan", "Generate a JSON file with all the responding hosts in a subnet")
	// app.Spec = ""
	app.Version("v version", appversion)

	app.Before = func() {

	}

	app.Command("scan", "Scan a subnet", cmdScan)
	app.Command("list", "List Saved Scans", cmdList)
	app.Command("compare", "Compare Scans", cmdCompare)
	app.Command("view", "View Details for a host", cmdView)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

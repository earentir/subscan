package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
)

func cmdView(cmd *cli.Cmd) {
	cmd.Spec = "SCANFILE TARGETIP"
	var (
		scanFile = cmd.StringArg("SCANFILE", "", "Scan file to view")
		targetIP = cmd.StringArg("TARGETIP", "", "Target IP to view")
	)

	cmd.Action = func() {
		output := readFiles([]string{*scanFile})
		result := findScanResultByIP(output[0], *targetIP)
		if result != nil {
			fmt.Println("DateTime:", result.DateTime)
			fmt.Println("HostIP:", result.HostIP)
			fmt.Println("OpenPorts:", result.OpenPorts)
			fmt.Println("PTRRecord:", result.PTRRecord)
			fmt.Println("ARecords:", result.ARecords)
			fmt.Println("CNAMEs:", result.CNAMEs)
			fmt.Println("IPMatch:", result.IPMatch)
			fmt.Println("TTL:", result.TTL)
		} else {
			fmt.Println("No results found for IP:", *targetIP)
		}
	}
}

func findScanResultByIP(output Output, targetIP string) *ScanResult {
	for _, result := range output.Results {
		if result.TargetIP == targetIP {
			return &result
		}
	}
	return nil
}

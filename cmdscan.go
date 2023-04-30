package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	cli "github.com/jawher/mow.cli"
)

func cmdScan(cmd *cli.Cmd) {

	cmd.Spec = "SUBNET [--dns] [--ports]"

	var (
		subnet      = cmd.StringArg("SUBNET", "", "Subnet to scan, example")
		dnsServer   = cmd.StringOpt("dns", "", "DNS Server to Use, example")
		portsToScan = cmd.StringsOpt("ports", []string{"22", "80", "443"}, "Ports to scan")
	)

	cmd.Action = func() {
		// fmt.Println("portsToScan:", *portsToScan)
		// os.Exit(1)

		localIP := getLocalIP()

		currentDate := time.Now().Format("2006-01-02")
		outputFilename := fmt.Sprintf("%s-%s-scan.json", strings.ReplaceAll(*subnet, "/", "-"), currentDate)

		ips, _ := getIPsFromSubnet(*subnet)
		results := make([]ScanResult, 0)

		var wg sync.WaitGroup
		for _, ip := range ips {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				result := ScanResult{
					DateTime: time.Now().Format(time.RFC3339),
					HostIP:   localIP,
					TargetIP: ip,
				}

				timeout := 1 * time.Second
				if isUp, ttl := tcpPing(ip, timeout); isUp {
					fmt.Printf("TTL for %s: %d ms\n", ip, ttl)
					result.OpenPorts = scanPorts(ip, *portsToScan)
					result.PTRRecord, result.ARecords, result.CNAMEs = getDNSRecords(ip, *dnsServer)
					result.IPMatch = checkIPMatch(ip, result.ARecords)
					result.TTL = ttl
					results = append(results, result)
				}
			}(ip)
		}

		wg.Wait()

		output := Output{
			ScanDate:  time.Now().Format(time.RFC3339),
			Subnet:    *subnet,
			HostIP:    localIP,
			DNSServer: *dnsServer,
			Results:   results,
		}

		naturalSort(&output)

		jsonData, _ := json.MarshalIndent(output, "", "  ")
		err := os.WriteFile(outputFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing JSON file:", err)
			return
		}

		fmt.Println("Scan completed, results saved to ", outputFilename)
	}
}

func naturalSort(output *Output) {
	sort.SliceStable(output.Results, func(i, j int) bool {
		a := strings.Split(output.Results[i].TargetIP, ".")
		b := strings.Split(output.Results[j].TargetIP, ".")
		for index := 0; index < len(a) && index < len(b); index++ {
			intA, _ := strconv.Atoi(a[index])
			intB, _ := strconv.Atoi(b[index])
			if intA != intB {
				return intA < intB
			}
		}
		return len(a) < len(b)
	})
}

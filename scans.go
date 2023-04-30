package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func printGroupedFiles(groupedFiles map[string][]string, dir string) {
	for key, files := range groupedFiles {
		fmt.Printf("IP-Subnet: %s\n", key)
		scanDates, hostIPs, dnsServers, resultCounts, err := processFiles(files, dir)
		if err != nil {
			fmt.Printf("Error processing files: %v\n", err)
			continue
		}

		fmt.Printf("ScanDates: %s\n", strings.Join(scanDates, ", "))
		for i := range files {
			fmt.Printf("Results count: %d\n", resultCounts[i])
			if dnsServers[i] == "" {
				fmt.Printf("HostIP: %s\n", hostIPs[i])
			} else {
				fmt.Printf("HostIP: %s, DNSServer: %s\n", hostIPs[i], dnsServers[i])
			}
		}
		fmt.Println()
	}
}

func processFiles(files []string, dir string) ([]string, []string, []string, []int, error) {
	var scanDates, hostIPs, dnsServers []string
	var resultCounts []int

	for _, file := range files {
		scanDate, hostIP, dnsServer, resultsCount, err := getMetadata(dir, file)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		scanDates = append(scanDates, scanDate)
		hostIPs = append(hostIPs, hostIP)
		dnsServers = append(dnsServers, dnsServer)
		resultCounts = append(resultCounts, resultsCount)
	}

	return scanDates, hostIPs, dnsServers, resultCounts, nil
}

func getMatchingFiles(dir string) ([]string, error) {
	var matchedFiles []string
	pattern := `^((?:\d{1,3}\.){3}\d{1,3})-(\d{1,2})-(\d{4})-(\d{2})-(\d{2})-scan\.json$`

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			matched, err := regexp.MatchString(pattern, d.Name())
			if err != nil {
				return err
			}

			if matched {
				matchedFiles = append(matchedFiles, d.Name())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matchedFiles, nil
}

func groupFilesByIPSubnet(files []string) map[string][]string {
	groupedFiles := make(map[string][]string)

	for _, file := range files {
		split := strings.Split(file, "-")
		key := fmt.Sprintf("%s-%s", split[0], split[1])

		if _, ok := groupedFiles[key]; !ok {
			groupedFiles[key] = []string{file}
		} else {
			groupedFiles[key] = append(groupedFiles[key], file)
		}
	}

	return groupedFiles
}

func getMetadata(dir, filename string) (string, string, string, int, error) {
	file, err := os.Open(filepath.Join(dir, filename))
	if err != nil {
		return "", "", "", 0, err
	}
	defer file.Close()

	var output Output
	err = json.NewDecoder(file).Decode(&output)
	if err != nil {
		return "", "", "", 0, err
	}

	return output.ScanDate, output.HostIP, output.DNSServer, len(output.Results), nil
}

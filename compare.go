package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

func compare(scans []string) {
	outputs := readFiles(scans)
	compareFiles(outputs, scans)
}

func readFiles(files []string) []Output {
	var outputs []Output
	for _, filename := range files {
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		var output Output
		err = json.Unmarshal(file, &output)
		if err != nil {
			panic(err)
		}
		output.HostIP = filepath.Base(filename)
		outputs = append(outputs, output)
	}
	return outputs
}

func compareFiles(outputs []Output, filenames []string) {
	// Loop over results in first file
	for _, res1 := range outputs[0].Results {
		// Check if TargetIP exists in all other files
		var res2s []*ScanResult
		for _, output := range outputs[1:] {
			var res2 *ScanResult
			for _, r := range output.Results {
				if r.TargetIP == res1.TargetIP {
					res2 = &r
					break
				}
			}
			if res2 == nil {
				fmt.Printf("TargetIP %s doesn't exist in file %s\n", res1.TargetIP, output.HostIP)
				continue
			}
			res2s = append(res2s, res2)
		}

		// Compare all results for TargetIP
		for i, res2 := range res2s {
			if !reflect.DeepEqual(res1.OpenPorts, res2.OpenPorts) {
				fmt.Printf("Open ports for TargetIP %s differ in file %s and %s\n", res1.TargetIP, outputs[i+1].HostIP, outputs[0].HostIP)
				fmt.Printf("%s: %v\n%s: %v\n", filenames[i+1], res2.OpenPorts, filenames[0], res1.OpenPorts)
			}
			if res1.PTRRecord != res2.PTRRecord {
				fmt.Printf("PTR record for TargetIP %s differs in file %s and %s\n", res1.TargetIP, outputs[i+1].HostIP, outputs[0].HostIP)
				fmt.Printf("%s: %s\n%s: %s\n", filenames[i+1], res2.PTRRecord, filenames[0], res1.PTRRecord)
			}
			if !reflect.DeepEqual(res1.ARecords, res2.ARecords) {
				fmt.Printf("A records for TargetIP %s differ in file %s and %s\n", res1.TargetIP, outputs[i+1].HostIP, outputs[0].HostIP)
				fmt.Printf("%s: %v\n%s: %v\n", filenames[i+1], res2.ARecords, filenames[0], res1.ARecords)
			}
		}
	}
}

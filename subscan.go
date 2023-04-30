package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/miekg/dns"
	"github.com/tatsushid/go-fastping"
)

type Output struct {
	ScanDate  string       `json:"scandate"`
	Subnet    string       `json:"subnet"`
	HostIP    string       `json:"hostip"`
	DNSServer string       `json:"dnsserver"`
	Results   []ScanResult `json:"results"`
}

type ScanResult struct {
	DateTime  string   `json:"datetime"`
	HostIP    string   `json:"hostip"`
	TargetIP  string   `json:"targetip"`
	OpenPorts []int    `json:"openports"`
	PTRRecord string   `json:"ptrrecord"`
	ARecords  []string `json:"arecords"`
	CNAMEs    []string `json:"cnames"`
	IPMatch   bool     `json:"ipmatch"`
	TTL       int      `json:"ttl"`
}

var (
	appversion string = "1.0.0"
	outputFile string
)

func main() {
	app := cli.App("subscan", "Generate a JSON file with all the responding hosts in a subnet")
	app.Spec = "SUBNET [-d | --dns] [-p | --ports]"
	app.Version("v version", appversion)

	subnet := app.StringArg("SUBNET", "", "Subnet to scan, example 192.168.178.0/24")
	dnsServer := app.StringOpt("d dns", "192.168.178.7:53", "DNS Server to Use, example 192.168.178.7:53")
	portsToScan := app.IntsOpt("p ports", []int{22, 80, 443}, "Ports to scan")

	app.Action = func() {
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
					result.OpenPorts = scanPorts(ip, []int{80, 443, 22})
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

		jsonData, _ := json.MarshalIndent(output, "", "  ")
		err := os.WriteFile(outputFilename, jsonData, 0644)
		if err != nil {
			fmt.Println("Error writing JSON file:", err)
			return
		}

		fmt.Println("Scan completed, results saved to ", outputFilename)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func getIPsFromSubnet(subnet string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	if len(ips) > 0 {
		ips = ips[1 : len(ips)-1]
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func tcpPing(ip string, timeout time.Duration) (bool, int) {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
		return false, 0
	}
	p.AddIPAddr(ra)

	var isUp bool
	var ttl int
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		isUp = true
		ttl = int(rtt / time.Millisecond)
	}

	err = p.Run()
	if err != nil {
		fmt.Println(err)
		return false, 0
	}

	return isUp, ttl
}

func scanPorts(ip string, ports []int) []int {
	var openPorts []int
	for _, port := range ports {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
		if err == nil {
			openPorts = append(openPorts, port)
			conn.Close()
		}
	}
	return openPorts
}

func getDNSRecords(ip, dnsServer string) (string, []string, []string) {
	var ptrRecord string
	var aRecords []string
	var cnames []string

	reverseIP, err := dns.ReverseAddr(ip)
	if err != nil {
		fmt.Println(err.Error())
		return "", nil, nil
	}
	ptrMsg := dns.Msg{}
	ptrMsg.SetQuestion(reverseIP, dns.TypePTR)
	ptrResponse, err := dns.Exchange(&ptrMsg, dnsServer)

	if err == nil {
		for _, ans := range ptrResponse.Answer {
			if ptr, ok := ans.(*dns.PTR); ok {
				ptrRecord = ptr.Ptr
				break
			}
		}
	}

	if ptrRecord != "" {
		aMsg := dns.Msg{}
		aMsg.SetQuestion(dns.Fqdn(ptrRecord), dns.TypeA)
		aResponse, err := dns.Exchange(&aMsg, dnsServer)

		if err == nil {
			for _, ans := range aResponse.Answer {
				if aRecord, ok := ans.(*dns.A); ok {
					aRecords = append(aRecords, aRecord.A.String())
				}
				if cname, ok := ans.(*dns.CNAME); ok {
					cnames = append(cnames, cname.Target)
				}
			}
		}
	}

	return ptrRecord, aRecords, cnames
}

func checkIPMatch(targetIP string, aRecords []string) bool {
	for _, aRecord := range aRecords {
		if targetIP == aRecord {
			return true
		}
	}
	return false
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "192.168.178.7:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

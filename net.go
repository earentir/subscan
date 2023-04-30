package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/tatsushid/go-fastping"
)

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

func parseCSVToInt(csvInput string) []int {
	var intSlice []int
	fields := strings.Split(csvInput, ",")
	for _, field := range fields {
		intField, err := strconv.Atoi(strings.TrimSpace(field))
		if err == nil {
			intSlice = append(intSlice, intField)
		}
	}
	return intSlice
}

func scanPorts(ip string, ports []string) []int {
	intPorts := parseCSVToInt(ports[0])
	var openPorts []int
	for _, port := range intPorts {
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

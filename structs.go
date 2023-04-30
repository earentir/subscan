package main

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

# subscan
Scan subnets and create a json with the results.


## examples
### Scan 192.168.178.0/24 and use the 192.168.178.7:53 DNS server to get dns records (PTR, A, CNAME etc)
```
./subscan scan 192.168.178.0/24 --dns 192.168.178.7:53

TTL for 192.168.178.180: 68 ms
TTL for 192.168.178.232: 63 ms
TTL for 192.168.178.7: 1 ms
TTL for 192.168.178.28: 1 ms
TTL for 192.168.178.229: 63 ms
TTL for 192.168.178.190: 68 ms
TTL for 192.168.178.3: 63 ms
TTL for 192.168.178.26: 1 ms
TTL for 192.168.178.238: 69 ms
TTL for 192.168.178.6: 64 ms
TTL for 192.168.178.30: 0 ms
TTL for 192.168.178.1: 1 ms
TTL for 192.168.178.27: 2 ms
TTL for 192.168.178.2: 1 ms
TTL for 192.168.178.179: 68 ms
TTL for 192.168.178.178: 69 ms
TTL for 192.168.178.5: 64 ms
TTL for 192.168.178.191: 67 ms
TTL for 192.168.178.181: 68 ms
TTL for 192.168.178.189: 69 ms
TTL for 192.168.178.112: 1 ms
TTL for 192.168.178.104: 1 ms
TTL for 192.168.178.233: 27 ms
TTL for 192.168.178.51: 58 ms
TTL for 192.168.178.57: 7 ms
TTL for 192.168.178.101: 1 ms
TTL for 192.168.178.173: 31 ms
TTL for 192.168.178.40: 7 ms
TTL for 192.168.178.167: 3 ms
TTL for 192.168.178.21: 3 ms
TTL for 192.168.178.129: 129 ms
TTL for 192.168.178.22: 0 ms
TTL for 192.168.178.52: 50 ms
TTL for 192.168.178.18: 4 ms
TTL for 192.168.178.193: 7 ms
TTL for 192.168.178.46: 0 ms
TTL for 192.168.178.55: 50 ms
TTL for 192.168.178.102: 0 ms
TTL for 192.168.178.192: 64 ms
Scan completed, results saved to  192.168.178.0-24-2023-04-30-scan.json
```

### List all scans
```
./subscan list

Found 4 scans

IP-Subnet: 192.168.178.0-24
ScanDates: 2023-04-29T10:46:01+03:00, 2023-04-30T15:46:01+03:00, 2023-04-30T17:49:38+03:00
Results count: 34
HostIP: 172.29.36.128
Results count: 33
HostIP: 172.29.36.128
Results count: 39
HostIP: 172.29.36.128, DNSServer: 192.168.178.7:53

IP-Subnet: 172.29.36.0-24
ScanDates: 2023-04-30T02:35:27+03:00
Results count: 1
HostIP: 172.29.36.128, DNSServer: 192.168.178.7:53
```

### Compare scans
```
./subscan compare 192.168.178.0-24-2023-04-30-scan.json 192.168.178.0-24-2023-04-29-scan.json

TargetIP 192.168.178.26 doesn't exist in file 192.168.178.0-24-2023-04-30-scan.json
Open ports for TargetIP 192.168.178.5 differ in file 192.168.178.0-24-2023-04-30-scan.json and 192.168.178.0-24-2023-04-29-scan.json
192.168.178.0-24-2023-04-30-scan.json: [22 443]
192.168.178.0-24-2023-04-29-scan.json: [22 80 443]
```

### View Target IP
```
./subscan view 192.168.178.0-24-2023-04-30-scan.json 192.168.178.1

DateTime: 2023-04-30T23:04:13+03:00
HostIP: 172.29.36.128
OpenPorts: [443 53]
PTRRecord: gw.ear.pm.
ARecords: [192.168.178.1]
CNAMEs: []
IPMatch: true
TTL: 0
```
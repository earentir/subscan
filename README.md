# subscan
Scan subnets and create a json with the results.


## examples
### Scan 192.168.178.0/24 and use the 192.168.178.7:53 DNS server to get dns records (PTR, A, CNAME etc)
```
./subscan scan 192.168.178.0/24 --dns 192.168.178.7:53
```

### List all scans
```
./subscan list
```

### Compare scans
```
./subscan compare 192.168.178.0-24-2023-04-30-scan.json 192.168.178.0-24-2023-04-29-scan.json
```

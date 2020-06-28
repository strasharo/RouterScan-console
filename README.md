It is a concole version of RouterScan for linux. It uses closed Router Core Library. Original soft is hosted here: [http://stascorp.com/load/1-1-0-56](http://stascorp.com/load/1-1-0-56)

## Usage

Scan single target

```bash
LD_LIBRARY_PATH=. ./routerscan scan --target <ip:port>
```

Scan multiple targets

```bash
LD_LIBRARY_PATH=. ./routerscan scan --input=./path/to/ip:port/list
# or
LD_LIBRARY_PATH=. echo "<ip:port>\n<ip2:port2>\n" | ./routerscan --input -
```

Show help

```bash
LD_LIBRARY_PATH=. ./routerscan scan -h
NAME:
   routerscan scan - run scanning routers in network

USAGE:
   routerscan scan [command options] [arguments...]

OPTIONS:
   --target value              <ip>:<port> (default: "192.168.1.1:80")
   --auth-basic value          ./path/to/<ip>\t<port>\n file with credentials dictionary for basic auth (default: "auth_basic.txt")
   --auth-digest value         ./path/to/<ip>\t<port>\nfile with credentials dictionary for digest auth (default: "auth_digest.txt")
   --auth-form value           ./path/to/<ip>\t<port>\n file with credentials dictionary for form auth (default: "auth_form.txt")
   --module-scanrouter         (default: true)
   --module-proxycheck         (default: false)
   --module-hnap               (default: true)
   --module-sqlite             (default: false)
   --module-hudson             (default: false)
   --module-phpmyadmin         (default: false)
   --st-enable-debug           (default: false)
   --st-debug-verbosity value  supported levels: 1, 2, 3 (default: 0)
   --st-user-agent value       (default: "Mozilla/5.0 (Windows NT 5.1; rv:9.0.1) Gecko/20100101 Firefox/9.0.1")
   --input value, -i value     /path/to/<ip>:<port>
 file with list of targets, or '-' for read targets from stdin (default: "-")
   --threads value, -t value  (default: 5)
   --help, -h                 show help (default: false)
```

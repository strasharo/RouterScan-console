It is a concole version of RouterScan for linux. It uses closed Router Core Library. Original soft is hosted here: [http://stascorp.com/load/1-1-0-56](http://stascorp.com/load/1-1-0-56)

## Usage

Scanning single target

```bash
    LD_LIBRARY_PATH=. ./routerscan scan --target <ip:port>
```

Scanning multiple targets

```bash
    LD_LIBRARY_PATH=. ./routerscan scan --input=./path/to/ip:port/list
    # or
    LD_LIBRARY_PATH=. echo "<ip:port>\n<ip2:port2>\n" | ./routerscan --input -
```
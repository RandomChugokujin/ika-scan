# ika-scan

ika-scan is a modern port scanner specifically designed to find open ports behind Squid proxies.
While this can be achieved with proxychains and Nmap, the performance is too slow to be acceptable.

![cover](img/cover.png)

## Build

```
$ git clone https://github.com/RandomChugokujin/ika-scan
$ cd ika-scan
$ go build
```

## Usage
```
$ ika-scan
Usage:
  ika-scan [OPTIONS]

Application Options:
  -u, --url=         The URL of the SQUID proxy (required)
  -w, --num-workers= Number of workers for port scanning, default is 100 (default: 100)
  -p, --num-ports=   Maximum number of ports scanned, default is 1000 (default: 1000)

Help Options:
  -h, --help         Show this help message
```

Example:
```
$ ika-scan -u http://10.10.108.208:3128 -p 10000
```

## Acknowledgement
This program built off of xct's [PoC](https://gist.github.com/xct/597d48456214b15108b2817660fdee00).

## License
This program is distributed with GPLv3 License.

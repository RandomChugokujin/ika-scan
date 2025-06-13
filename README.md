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
$ ika-scan -h
Usage:
  ika-scan [OPTIONS]

Application Options:
  -u, --url=         The URL of the SQUID proxy (required)
  -p, --ports=       A list of ports to scan, support both commas and ranges with -
                     (required).
  -w, --num-workers= Number of workers for port scanning. (default: 100)

Help Options:
  -h, --help         Show this help message

```

Example:
```
$ ika-scan -u http://10.10.108.208:3128 -p 22,80,443,9000-10000
```

## Acknowledgement
This program built off of xct's [PoC](https://gist.github.com/xct/597d48456214b15108b2817660fdee00).

## License
This program is distributed with GPLv3 License.

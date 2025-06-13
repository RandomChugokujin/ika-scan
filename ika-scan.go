package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Url 	   	string `short:"u" long:"url" description:"The URL of the SQUID proxy (required)" required:"true"`
	Ports      	string `short:"p" long:"ports" description:"A list of ports to scan, support both commas and ranges with - (required)." required:"true"`
	NumWorkers 	int    `short:"w" long:"num-workers" default:"100" description:"Number of workers for port scanning, default is 100."`
	// Verbose  	bool   `short:"v" long:"verbose" default:"False" description:"Enable verbose output."`
}

func port_parse(port_arg string) ([]uint16, error) {
	port_slice := make([]uint16, 10)
	var err error

	ports_split_by_comma := strings.Split(port_arg, ",")
	for _, port := range ports_split_by_comma {
		if strings.Contains(port, "-"){
			ports_split_by_dash := strings.Split(port, "-")
			if len(ports_split_by_dash) != 2 {
				err = errors.New("Invalid Port Range")
				return port_slice, err
			}
			range_start, err_start := strconv.ParseUint(ports_split_by_dash[0], 10, 16)
			range_end, err_end := strconv.ParseUint(ports_split_by_dash[1], 10, 16)
			if err_start != nil || err_end != nil || range_start >= range_end {
				err := errors.New("Invalid Port Range")
				return port_slice, err
			}

			for i := range_start; i <= range_end; i++ {
				port_slice = append(port_slice, uint16(i))
			}
		} else {
			// Port Check correct size
			port_num, err := strconv.ParseUint(port, 10, 16)
			if err != nil{
				fmt.Printf("Invalid Port Number %v\n", port)
				return port_slice, err
			}
			port_slice = append(port_slice, uint16(port_num))
		}
	}
	return port_slice, nil
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	banner  := `
_________ _        _______         _______  _______  _______  _
\__   __/| \    /\(  ___  )       (  ____ \(  ____ \(  ___  )( (    /|
   ) (   |  \  / /| (   ) |       | (    \/| (    \/| (   ) ||  \  ( |
   | |   |  (_/ / | (___) | _____ | (_____ | |      | (___) ||   \ | |
   | |   |   _ (  |  ___  |(_____)(_____  )| |      |  ___  || (\ \) |
   | |   |  ( \ \ | (   ) |             ) || |      | (   ) || | \   |
___) (___|  /  \ \| )   ( |       /\____) || (____/\| )   ( || )  \  |
\_______/|_/    \/|/     \|       \_______)(_______/|/     \||/    )_)

			Gas Gas Gas ~~~
	`
	fmt.Println(banner)

	proxyURL, err := url.Parse(opts.Url)

	if err != nil {
		fmt.Printf("Failed to parse proxy URL: %v\n", err)
		return
	}

	// Parse Ports
	ports_to_scan, err := port_parse(opts.Ports)
	if err != nil {
		fmt.Printf("Failure to parse ports provided: %v\n", err)
		return
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: transport}
	openPorts := make([]uint16, 0)

	bar := pb.StartNew(len(ports_to_scan))
	sem := make(chan struct{}, opts.NumWorkers)
	var wg sync.WaitGroup

	for _, port := range ports_to_scan {
		wg.Add(1)
		go func(p uint16) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() {
				<-sem
				bar.Increment()
			}()

			address := fmt.Sprintf("127.0.0.1:%d", p)
			r, err := client.Get(fmt.Sprintf("http://%s", address))
			if err != nil {
				return
			}
			data, _ := io.ReadAll(r.Body)
			dataStr := string(data)
			if strings.Contains(dataStr,"The requested URL could not be retrieved"){
				return;
			}
			defer r.Body.Close()
			openPorts = append(openPorts, p)
			// if opts.Verbose{
			// 	fmt.Printf("Port %d found!\n", p)
			// }
		}(port)
	}
	wg.Wait()
	bar.Finish()

	fmt.Println("Open ports:")
	for _, port := range openPorts {
		fmt.Println(port)
	}
}

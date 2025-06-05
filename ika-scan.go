package main

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
	"strings"
	"io"
	"os"
	"github.com/cheggaaa/pb/v3"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Url        string `short:"u" long:"url" description:"The URL of the SQUID proxy (required)" required:"true"`
	NumWorkers int    `short:"w" long:"num-workers" default:"100" description:"Number of workers for port scanning, default is 100"`
	NumPorts   int    `short:"p" long:"num-ports" default:"1000" description:"Maximum number of ports scanned, default is 1000"`
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


	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 3 * time.Second,
		}).DialContext,
	}
	client := &http.Client{Transport: transport}
	openPorts := make([]int, 0)

	bar := pb.StartNew(opts.NumPorts)
	sem := make(chan struct{}, opts.NumWorkers)
	var wg sync.WaitGroup
	for port := 1; port <= opts.NumPorts; port++ {
		wg.Add(1)
		go func(p int) {
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
			fmt.Printf("Port %d found!\n", p)

		}(port)
	}
	wg.Wait()
	bar.Finish()

	fmt.Println("Open ports:")
	for _, port := range openPorts {
		fmt.Println(port)
	}
}

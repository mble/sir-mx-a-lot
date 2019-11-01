package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mble/sir-mx-a-lot/pkg/resolver"
)

func ingest(done <-chan bool) (<-chan string, error) {
	flag.Parse()
	var data []byte
	var err error
	domains := make(chan string)
	switch flag.NArg() {
	case 0:
		data, err = ioutil.ReadAll(os.Stdin)
		break
	case 1:
		data, err = ioutil.ReadFile(flag.Arg(0))
		break
	default:
		err = fmt.Errorf("input must be from stdin or file")
		break
	}
	go func() {
		for _, domain := range strings.Split(string(data), "\n") {
			select {
			case <-done:
				return
			case domains <- domain:
			}
		}
		close(domains)
	}()
	return domains, err
}

func lookMX(done <-chan bool, domains <-chan string, resolver *net.Resolver, workerID int) <-chan string {
	const timeout = 500 * time.Millisecond
	processed := make(chan string)
	go func() {
		for domain := range domains {
			select {
			case <-done:
				return
			case processed <- domain:
				if domain != "" {
					mx(domain, resolver, timeout)
				}
			}
		}
		close(processed)
	}()
	return processed
}

func merge(done <-chan bool, channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	processedDomains := make(chan string)
	multiplex := func(c <-chan string) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case processedDomains <- i:
			}
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(processedDomains)
	}()
	return processedDomains
}

func mx(domain string, resolver *net.Resolver, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mxs, err := resolver.LookupMX(ctx, domain)
	if err != nil {
		fmt.Printf("%s: %s\n", domain, err)
	}
	for _, mx := range mxs {
		fmt.Printf("%s: %s\n", domain, mx.Host)
	}
}

func main() {
	var numWorkers = runtime.NumCPU()

	done := make(chan bool)
	defer close(done)

	resolver := resolver.NewResolver()

	domains, err := ingest(done)
	if err != nil {
		panic(err)
	}

	workers := make([]<-chan string, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = lookMX(done, domains, resolver, i)
	}

	for range merge(done, workers...) {
	}
}

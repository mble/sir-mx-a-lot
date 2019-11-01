package pattern

import (
	"context"
	"fmt"
	"net"
	"time"
)

// LookMX receives a channel of domains and ranges over them, looking up the MX records for each domain and returning
// a channel of results
func LookMX(done <-chan bool, domains <-chan string, resolver *net.Resolver, workerID int) <-chan string {
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

package pipeline

import (
	"context"
	"net"
	"time"

	"github.com/mble/sir-mx-a-lot/internal/logger"
)

// LookMX receives a channel of domains and ranges over them, looking up the MX records for each domain and returning
// a channel of results
func LookMX(done <-chan bool, domains <-chan Domain, resolver *net.Resolver, workerID int) <-chan Domain {
	const timeout = 500 * time.Millisecond
	processed := make(chan Domain)
	go func() {
		for domain := range domains {
			select {
			case <-done:
				return
			case processed <- domain:
				if domain.Name != "" {
					mx(domain, resolver, timeout)
				}
			}
		}
		close(processed)
	}()
	return processed
}

func mx(domain Domain, resolver *net.Resolver, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	mxs, err := resolver.LookupMX(ctx, domain.Name)
	if err != nil {
		domain.Error = err.Error()
	}
	for _, mx := range mxs {
		domain.MXHosts = append(domain.MXHosts, mx.Host)
	}
	logger.Log(domain)
}

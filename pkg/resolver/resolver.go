package resolver

import (
	"net"

	"github.com/benburkert/dns"
)

// NewResolver returns a new custom resolver using Google's and Cloudflares DNS resolvers.
func NewResolver() *net.Resolver {
	return &net.Resolver{
		PreferGo: true,

		Dial: (&dns.Client{
			Transport: &dns.Transport{
				Proxy: dns.NameServers{
					&net.UDPAddr{IP: net.IPv4(8, 8, 8, 8), Port: 53}, // Google primary
					&net.UDPAddr{IP: net.IPv4(8, 8, 4, 4), Port: 53}, // Google backup
					&net.UDPAddr{IP: net.IPv4(1, 1, 1, 1), Port: 53}, // Cloudflare primary
					&net.UDPAddr{IP: net.IPv4(1, 0, 0, 1), Port: 53}, // Cloudflare backup
				}.RoundRobin(),
			},
		}).Dial,
	}
}

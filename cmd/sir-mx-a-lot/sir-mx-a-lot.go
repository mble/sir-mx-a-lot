package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	"github.com/mble/sir-mx-a-lot/pkg/resolver"
)

func ingest() ([]string, error) {
	flag.Parse()
	var data []byte
	var err error
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
	return strings.Split(string(data), "\n"), err
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
	const timeout = 500 * time.Millisecond

	resolver := resolver.NewResolver()

	data, err := ingest()
	if err != nil {
		panic(err)
	}

	for _, domain := range data {
		if domain != "" {
			mx(domain, resolver, timeout)
		}
	}
}

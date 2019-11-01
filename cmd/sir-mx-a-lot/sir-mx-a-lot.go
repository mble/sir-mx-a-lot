package main

import (
	"context"
	"fmt"

	"github.com/mble/sir-mx-a-lot/pkg/resolver"
)

func main() {
	resolver := resolver.NewResolver()
	mxs, err := resolver.LookupMX(context.Background(), "blwt.io")
	if err != nil {
		panic(err)
	} else {
		for _, mx := range mxs {
			fmt.Printf("%s\n", mx.Host)
		}
	}
}

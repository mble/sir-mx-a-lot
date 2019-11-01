package main

import (
	"runtime"

	"github.com/mble/sir-mx-a-lot/internal/pattern"
	"github.com/mble/sir-mx-a-lot/pkg/resolver"
)

func main() {
	var numWorkers = runtime.NumCPU()

	done := make(chan bool)
	defer close(done)

	resolver := resolver.NewResolver()

	domains, err := pattern.Ingest(done)
	if err != nil {
		panic(err)
	}

	workers := make([]<-chan string, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = pattern.LookMX(done, domains, resolver, i)
	}

	for range pattern.Merge(done, workers...) {
	}
}

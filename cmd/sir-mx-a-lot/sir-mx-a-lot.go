package main

import (
	"runtime"

	"github.com/mble/sir-mx-a-lot/internal/pipeline"
	"github.com/mble/sir-mx-a-lot/pkg/resolver"
)

func main() {
	var numWorkers = runtime.NumCPU()

	done := make(chan bool)
	defer close(done)

	resolver := resolver.NewResolver()

	domains, err := pipeline.Ingest(done)
	if err != nil {
		panic(err)
	}

	workers := make([]<-chan pipeline.Domain, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = pipeline.LookMX(done, domains, resolver, i)
	}

	for range pipeline.Merge(done, workers...) {
	}
}

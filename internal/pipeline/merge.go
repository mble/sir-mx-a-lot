package pipeline

import "sync"

// Merge multiplexes incoming channels, passing work across them and returning
// a single channel of performed work
func Merge(done <-chan bool, channels ...<-chan Domain) <-chan Domain {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	processedDomains := make(chan Domain)
	multiplex := func(c <-chan Domain) {
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

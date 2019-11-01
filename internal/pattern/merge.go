package pattern

import "sync"

// Merge multiplexes incoming channels, passing work across them and returning
// a single channel of performed work
func Merge(done <-chan bool, channels ...<-chan string) <-chan string {
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

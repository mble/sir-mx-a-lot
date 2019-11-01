package pattern

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Ingest reads domains from stdin or a file, returning
// a string channel and error
func Ingest(done <-chan bool) (<-chan string, error) {
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

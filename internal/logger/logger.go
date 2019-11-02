package logger

import (
	"encoding/json"
	"fmt"
)

// Log marshals a domain struct into JSON and logs it out.
func Log(domain interface{}) {
	data, _ := json.Marshal(domain)
	fmt.Printf("%s\n", data)
}

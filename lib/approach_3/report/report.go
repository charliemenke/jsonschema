package report

import (
	"fmt"
)

type BadKeyVal struct {
	Key   string
	Value any
	Error error
}

func StdOutReporter(e BadKeyVal) error {
	fmt.Printf("\nfailed to validate key <%s> with value <%v>: %s\n", e.Key, e.Value, e.Error)
	return nil
}

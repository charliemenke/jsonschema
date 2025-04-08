package approach_3

import (
	"cmenke/go-playground/lib/approach_3/report"
	"cmenke/go-playground/lib/approach_3/schema"
	"fmt"
)

type Listing struct {
	DocId string   `json:"docid"`
	Mls   string   `json:"mls"`
	Data  []KeyVal `json:"data"`
}

type KeyVal struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}


type Reporter interface {
	Report(e report.BadKeyVal) error
}

type ReporterFunc func(e report.BadKeyVal) error

func (rf ReporterFunc) Report(e report.BadKeyVal) error {
	return rf(e)
}

func (kv *KeyVal) Validate(s *schema.Schema, r Reporter, coerce bool) (KeyVal, error) {
	// check if key is specified in passed schema
	s, ok := (*s.Properties)[kv.Key]
	if !ok {
		// if no schema specified, we accept the key:value as is
		fmt.Printf("\nno schema mapping for key <%s>, continuing\n", kv.Key)
		return *kv, nil
	}

	// if schema does exist for key:value, evaluate the value
	// recursivly via the passed schema
	val, err := s.Eval(kv.Value, coerce)
	if err != nil {
		// if schema for property fails to eval, report error to reporter
		// and return the error
		r.Report(report.BadKeyVal{
			Key:   kv.Key,
			Value: kv.Value,
			Error: err,
		})
		return *kv, err
	}

    // else if key:value is valid, return a new key:value as the Value
    // member could be coerced
	return KeyVal{
		Key:   kv.Key,
		Value: val,
	}, nil
}

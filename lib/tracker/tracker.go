package tracker

import "fmt"

type BadKeyVals struct {
	Key   string
	Value any
	Msg   error
}

type StdOutTracker struct {
	InvalidKeyVals []BadKeyVals
}

func (t *StdOutTracker) TrackInvalidType(key string, value any, msg error) error {
	t.InvalidKeyVals = append(t.InvalidKeyVals, BadKeyVals{
		Key:   key,
		Value: value,
		Msg:   msg,
	})
	return nil
}

func (t *StdOutTracker) Output() error {
	for _, e := range t.InvalidKeyVals {
		fmt.Printf("unable to validate key <%s> with value <%v> as a float64: %s\n", e.Key, e.Value, e.Msg)
	}
	return nil
}

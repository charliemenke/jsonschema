package approach_1

import (
	"cmenke/go-playground/lib/tracker"
	"encoding/json"
	"fmt"
	"strconv"
)

type InvalidTypeTracker interface {
	TrackInvalidType(key string, val any, msg error) error
	Output() error
}

type Listing struct {
	DocId string    `json:"docid"`
	Mls   string    `json:"mls"`
	Data  []KeyVals `json:"data"`
}

func (l *Listing) ValidFields() []KeyVals {
	validFields := make([]KeyVals, 0, len(l.Data))
	for _, e := range l.Data {
		if e.Key != "REMOVE_ME" {
			validFields = append(validFields, e)
		}
	}
	return validFields
}

type KeyVals struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (r *KeyVals) SetForRemoval() {
	r.Key = "REMOVE_ME"
	r.Value = nil
}

func (r *KeyVals) BubbleUpError(tracker InvalidTypeTracker, key string, val any, msg error) error {
	err := tracker.TrackInvalidType(key, val, msg)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *KeyVals) UnmarshalJSON(b []byte) error {
	// fmt.Printf("attempting to unmarshal: %s\n", string(b))
	// create alias of type to prevent recursively calling this unmarshal method and get a stack overflow
	type KeyValsAlias KeyVals
	var aux KeyValsAlias

	// unmarshal key:value element into temp alias struct
	err := json.Unmarshal(b, &aux)
	if err != nil {
		fmt.Printf("error unmarshaling: %s", err)
		return err
	}

	// attempt to validate type
	switch aux.Key {
	case "ListPrice", "OriginalListPrice":
		// if already float value, append and continue
		if _, ok := aux.Value.(float64); ok {
			r = (*KeyVals)(&aux)
			return nil
		}

		// if string, try to convert to float, append and continue
		if strV, ok := aux.Value.(string); ok {
			intVal, err := strconv.Atoi(strV)
			if err != nil {
				r.BubbleUpError(&TypeTracker, aux.Key, aux.Value, err)
				r.SetForRemoval()
				return nil
			}
			r.Key = aux.Key
			r.Value = intVal
			return nil
		}

		r.SetForRemoval()
		return nil
	case "Appliances":
		// if slice already, just set value
		if _, ok := aux.Value.([]string); ok {
			r = (*KeyVals)(&aux)
			return nil
		}

		// check if type is string
		if strVal, ok := aux.Value.(string); ok {
			// if string, attempt to read json stringified array
			var arr []string
			if err := json.Unmarshal([]byte(strVal), &arr); err != nil {
				r.BubbleUpError(&TypeTracker, aux.Key, aux.Value, err)
				r.SetForRemoval()
				return nil
			}
			r.Key = aux.Key
			r.Value = arr
			return nil
		}

		r.SetForRemoval()
		return nil
	default:
		// warn we have an unmapped key
		fmt.Printf("no explicit mapping found for key: %s with value: %+v\n", aux.Key, aux.Value)
		r = (*KeyVals)(&aux)
	}

	return nil
}

var TypeTracker tracker.StdOutTracker = tracker.StdOutTracker{}

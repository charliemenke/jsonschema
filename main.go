package main

import (
	"encoding/json"
	"fmt"
	"os"

	"cmenke/go-playground/lib/approach_1"
	"cmenke/go-playground/lib/approach_3"
	"cmenke/go-playground/lib/approach_3/report"
	"cmenke/go-playground/lib/approach_3/schema"
	// "cmenke/go-playground/lib/approach_2"
	// "cmenke/go-playground/lib/tracker"
)

func main() {
	testData := []byte(`{
        "mls": "rets-properties-test",
        "docid": "1234",
        "data": [
            {
                "key": "ListPrice",
                "value": "100000"
            },
            {
                "key": "ListPrice",
                "value": 100000.10
            },
            {
                "key": "ListPrice",
                "value": 100000
            },
            {
                "key": "Appliances",
                "value": "[\"a\", \"b\", \"c\"]"
            },
            {
                "key": "Appliances",
                "value": ["a", 10, "c"]
            },
            {
                "key": "Appliances",
                "value": ["a", "b", "c"]
            },
            {
                "key": "DontMapMe",
                "value": "555-555-5555"
            },
            {
                "key": "Geo",
                "value": [
					{
						"Timezone": "im a string now, screw you",
						"identifier": "GeoNSRF"
					}
				]	
            },
            {
                "key": "Geo",
                "value": [
					{
						"Timezone": {
							"TimezoneCode": "P",
							"TimezoneStdOffset": "-8",
							"Name": "America/Los_Angeles",
							"ObservesDLS": true
						},
						"identifier": "GeoNSRF"
					}
				]	
            }
        ]
    }`)

	fmt.Println("------Approach One--------")
	var l approach_1.Listing
	err := json.Unmarshal(testData, &l)
	if err != nil {
		fmt.Println("approach_1: error unmarshaling test listing")
		panic(err)
	}
	fmt.Printf("before clean: %+v\n", l)
	validFields := l.ValidFields()
	if len(validFields) > 0 {
		l.Data = validFields
	}
	fmt.Printf("after clean: %+v\n", l)
	_ = approach_1.TypeTracker.Output()

	// fmt.Println("\n------Approach two--------")
	// var l2 approach_2.Listing
	// err = json.Unmarshal(testData, &l2)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("before clean: %+v\n", l2)
	// typeTracker := tracker.StdOutTracker{}
	// _ = l2.ValidateTypes(&typeTracker)
	// fmt.Printf("after clean: %+v\n", l2)
	// _ = typeTracker.Output()

	fmt.Println("\n------Approach three--------")
	// read in json string listing
	var l3 approach_3.Listing
	err = json.Unmarshal(testData, &l3)
	if err != nil {
		fmt.Println("approach_3: error unmarshaling test listing")
		panic(err)
	}
	// read in schema
	var mapping schema.Schema
	file, err := os.ReadFile("./schema.json")
	if err != nil {
		fmt.Println("approach_3: error reading schema file")
		panic(err)
	}
	err = json.Unmarshal(file, &mapping)
	if err != nil {
		fmt.Println("approach_3: error unmarshaling schema")
		panic(err)
	}

	// debug print listing before validation
	str, _ := json.MarshalIndent(l3, "", "\t") 
	fmt.Printf("before clean:\n%s\n", str)

	// loop over each key:value, validating each one
	// if key:value is "invalid", remove it
	newData := make([]approach_3.KeyVal, 0, len(l3.Data))
	badData := []approach_3.KeyVal{}
	for _, e := range l3.Data {
		validatedKeyVal, err := e.Validate(&mapping, approach_3.ReporterFunc(report.StdOutReporter))
		if err != nil {
			badData = append(badData, e)
			continue
		}
		newData = append(newData, validatedKeyVal)
	}

	// assign validated data
	l3.Data = newData
	str, _ = json.MarshalIndent(l3, "", "\t") 
	fmt.Printf("after clean:\n%s\n", str)
}

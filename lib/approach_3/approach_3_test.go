package approach_3

import (
	"cmenke/go-playground/lib/approach_3/report"
	"cmenke/go-playground/lib/approach_3/schema"
	"encoding/json"
	"reflect"
	"testing"
)

func TestApproach3(t *testing.T) {
	testCases := []struct {
		description    string
		input          []KeyVal
		schema         []byte
		expectedOutput []KeyVal
	}{
		//////////////////////////////////////////////////////////
		//  simple value type to schema type acceptance testing //
		//////////////////////////////////////////////////////////
		{
			description: "it should accept a string value as a string type",
			input: []KeyVal{
				{
					Key:   "OfficeName",
					Value: "Foo",
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"OfficeName": { "type": "string" }
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "OfficeName",
					Value: "Foo",
				},
			},
		},
		{
			description: "it should accept an interger value as an integer type",
			input: []KeyVal{
				{
					Key:   "ListPrice",
					Value: 100000,
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"ListPrice": { "type": "integer" }
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "ListPrice",
					Value: 100000,
				},
			},
		},
		{
			description: "it should accept a float value as a number type",
			input: []KeyVal{
				{
					Key:   "ListPrice",
					Value: 100000.10,
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"ListPrice": { "type": "number" }
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "ListPrice",
					Value: 100000.10,
				},
			},
		},
		{
			description: "it should accept a boolean value as a boolean type",
			input: []KeyVal{
				{
					Key:   "DisplayYN",
					Value: false,
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"DisplayYN": { "type": "boolean" }
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "DisplayYN",
					Value: false,
				},
			},
		},
		{
			description: "it should accept an array of strings value as an array of strings type",
			input: []KeyVal{
				{
					Key:   "Appliances",
					Value: []any{"microwave", "oven", "fridge"},
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"Appliances": {
						"type": "array",
						"items": { "type": "string" }
					}
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "Appliances",
					Value: []any{"microwave", "oven", "fridge"},
				},
			},
		},
		{
			description: "it should accept an object value object type",
			input: []KeyVal{
				{
					Key:   "someObj",
					Value: map[string]any{
						"someKey": "someVal",
					},
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"someObj": {
						"type": "object",
						"properties": {
							"someKey": { "type": "string"}
						}
					}
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "someObj",
					Value: map[string]any{
						"someKey": "someVal",
					},
				},
			},
		},
		//////////////////////////////////////////////////////////
		//  simple value type to schema type rejection testing  //
		//////////////////////////////////////////////////////////
		{
			description: "it should reject a number value as a string type",
			input: []KeyVal{
				{
					Key:   "OfficeName",
					Value: 0,
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"OfficeName": { "type": "string" }
    			}
			}`),
			expectedOutput: []KeyVal{},
		},
		{
			description: "it should reject a float value as a integer type",
			input: []KeyVal{
				{
					Key:   "ListPrice",
					Value: 100000.10,
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"ListPrice": { "type": "integer" }
    			}
			}`),
			expectedOutput: []KeyVal{},
		},
		{
			description: "it should reject a string value as a boolean type",
			input: []KeyVal{
				{
					Key:   "DisplayYN",
					Value: "Nah",
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"DisplayYN": { "type": "boolean" }
    			}
			}`),
			expectedOutput: []KeyVal{},
		},
		{
			description: "it should reject a string value as an array of strings type",
			input: []KeyVal{
				{
					Key:   "Appliances",
					Value: "I should be an array :)",
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"Appliances": {
						"type": "array",
						"items": { "type": "string" }
					}
    			}
			}`),
			expectedOutput: []KeyVal{},
		},
		{
			description: "it should reject an array value as an object type",
			input: []KeyVal{
				{
					Key:   "someObj",
					Value: []any{"I", "should", "be", "an", "obj"},
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"someObj": {
						"type": "object",
						"properties": {
							"someKey": { "type": "string"}
						}
					}
    			}
			}`),
			expectedOutput: []KeyVal{},
		},
		//////////////////////////////////////////////////////////
		//                no schema spec testing                //
		//////////////////////////////////////////////////////////
		{
			description: "it should accept any key value not specified in schema",
			input: []KeyVal{
				{
					Key:   "imNotMapped",
					Value: "so add me!",
				},
			},
			schema: []byte(`{
 			   "properties": {}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "imNotMapped",
					Value: "so add me!",
				},
			},
		},
		{
			description: "it should accept any key value in object if its properties are not specified in schema",
			input: []KeyVal{
				{
					Key:   "someObj",
					Value: map[string]any{
						"imNotMapped": "so im allowed in",
					},
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"someObj": {
						"type": "object"
					}
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "someObj",
					Value: map[string]any{
						"imNotMapped": "so im allowed in",
					},
				},
			},
		},
		{
			description: "it should accept any element in array if its items are not specified in schema",
			input: []KeyVal{
				{
					Key:   "someArr",
					Value: []any{"we", 1, "are", false, []string{"allowed"}},
				},
			},
			schema: []byte(`{
 			   "properties": {
        			"someArr": {
						"type": "array"
					}
    			}
			}`),
			expectedOutput: []KeyVal{
				{
					Key:   "someArr",
					Value: []any{"we", 1, "are", false, []string{"allowed"}},
				},
			},
		},
	}

	for _, testCase := range testCases {
		var s schema.Schema
		err := json.Unmarshal(testCase.schema, &s)
		if err != nil {
			t.Fatalf("failed to unmarshal schema: %s", err)
		}

		testOutput := []KeyVal{}
		for _, kv := range testCase.input {
			validatedKV, err := kv.Validate(&s, ReporterFunc(report.StdOutReporter))
			if err != nil {
				continue
			}
			testOutput = append(testOutput, validatedKV)
		}

		if !reflect.DeepEqual(testOutput, testCase.expectedOutput) {
			t.Fatalf("testOutput <%v> does not match expected output <%v>", testOutput, testCase.expectedOutput)
		}
	}
}

package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type SchemaMap map[string]*Schema

type Type []string

func (t *Type) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*t = Type{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err == nil {
		*t = Type(multi)
		return nil
	}

	return errors.New(fmt.Sprintf("type field not single string or array of strings, value <%v>", string(data)))
}

type Schema struct {
	Properties *map[string]*Schema `json:"properties,omitempty"`
	Type       Type                `json:"type,omitempty"`
	Items      *Schema             `json:"items,omitempty"`
}

func coerceType(v string, toType Type) (any, error) {
	// loop over toType (because it can be a slice).
	// return first coerce that works
	var err error
	for _, t := range toType {
		switch t {
		case "null": // if we want value to be null, just set to nil
			return nil, nil
		case "boolean": // only parse stringifed true and false
			if v == "true" {
				return true, nil
			} else if v == "false" {
				return false, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to boolean", v)), err)
		case "integer": // try to parse integer first as it is more specific than number
			intVal, intCoerceError := strconv.Atoi(v)
			if intCoerceError == nil {
				return intVal, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to integer: %s", v, intCoerceError)), err)
		case "number":
			intVal, intCoerceError := strconv.Atoi(v)
			if intCoerceError == nil {
				return intVal, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to integer: %s", v, intCoerceError)), err)
			floatVal, floatCoerceError := strconv.ParseFloat(v, 64)
			if floatCoerceError == nil { // try parse float
				return floatVal, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to float: %s", v, floatCoerceError)), err)
		case "array":
			var arr []any
			arrCoerceError := json.Unmarshal([]byte(v), &arr)
			if arrCoerceError != nil {
				return arr, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to array: %s", v, arrCoerceError)), err)
		case "object":
			var obj map[string]any
			objCoerceError := json.Unmarshal([]byte(v), &obj)
			if objCoerceError != nil {
				return obj, nil
			}
			err = errors.Join(errors.New(fmt.Sprintf("failed to parse value <%s> to obj: %s", v, objCoerceError)), err)
		default:
			err = errors.Join(errors.New(fmt.Sprintf("cannot coerce value <%s> to unknown type <%s>", v, t)), err)
		}
	}	
	return nil, err
}

func (s *Schema) Eval(v any, coerce bool) (any, error) {
	var err error
	// handle base type check
	if s.Type != nil {
		valType, err := evalType(s, v)
		if err != nil {
			// if we do not want to coerce OR no valType was returned
			// OR valType is not a string return original eval error
			// since we will not try to coerce a non string type
			if !coerce || valType == "" || valType != "string" {
				return v, err
			}

			// otherwise, attempt to cast string to intended type
			// if we fail to coerce value, return original error
			// joined with coerce error
			var coerceErr error
			v, coerceErr = coerceType(v.(string), s.Type)	
			if coerceErr != nil {
				return v, errors.Join(err, coerceErr)
			}

			// if we do coerceType, do nothing since `v` is now
			// correct type
		}
	}
	// handle array
	if s.Items != nil {
		v, err = evalArray(s, v)
		if err != nil {
			return nil, err
		}
	}
	// handle properties
	if s.Properties != nil {
		v, err = evalObject(s, v)
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func evalType(s *Schema, val any) (string, error) {
	if s == nil {
		return "", errors.New("schema is nil, cannot eval type")
	}

	// if no type specified, accept it as is
	if len(s.Type) == 0 {
		fmt.Printf("no type specified\n")
		return "", nil
	}

	// get underlying type of val and check against specified
	// schema types
	valType := GetDataType(val)
	for _, t := range s.Type {
		if valType == t {
			return valType, nil
		}
	}

	// if no matches, return error specifying type mismatch
	return valType, errors.New(fmt.Sprintf("type mismatch, the value <%v> has the type <%s> which does not match expected type(s) <%v>", val, valType, s.Type))
}

func evalObject(s *Schema, val any) (any, error) {
	if s == nil {
		return nil, errors.New("schema is nil, cannot eval object")
	}

	var err error

	// ensure val is of object type
	valObj, ok := val.(map[string]any)
	if !ok {
		return nil, errors.New(fmt.Sprintf("value <%v> cannot be evaluated as an array", val))
	}

	// if properties is specified by schema, evaluate them
	if s.Properties != nil {
		val, err = evalProperties(s, valObj)
		if err != nil {
			// if unable to evaluate all properties, return error
			return val, errors.Join(fmt.Errorf("failed to validate object: %v", valObj), err)
		}
	}

	// else return val if no properties are specified meaning the obj type can contain anything
	return val, nil

}

func evalProperties(s *Schema, val map[string]any) (any, error) {
	if s == nil {
		return nil, errors.New("\tschema is nil, cannot eval properties")
	}

	if s.Properties == nil {
		fmt.Printf("no object schema specified\n")
		return val, nil
	}

	// for each key:value in object
	validKeyVals := map[string]any{}
	errs := []error{}
	for objK, objV := range val {
		if _, ok := (*s.Properties)[objK]; !ok {
			// if obj key not specified in properties schema, add to validKeyVals and continue
			// as no schema was specified
			validKeyVals[objK] = objV
			continue
		}
		// otherwise, attempt to evaluate the key:value as per schema spec
		v, err := (*s.Properties)[objK].Eval(objV)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		validKeyVals[objK] = v
	}

	if len(errs) > 0 {
		return nil, errors.New(fmt.Sprintf("\tcould not validate all key:vals in obj: %v", errs))
	}
	return validKeyVals, nil
}

func evalArray(s *Schema, val any) (any, error) {
	if s == nil {
		return nil, errors.New("schema is nil, cannot eval array")
	}

	var err error
	// ensure val is array type
	valItems, ok := val.([]any)
	if !ok {
		return nil, errors.New(fmt.Sprintf("value <%v> with type <%s> cannot be evaluated as an array", val, reflect.TypeOf(val)))
	}

	// enture items are correct schema if schema specifies one
	if s.Items != nil {
		val, err = evalItems(s, valItems)
		if err != nil {
			return val, errors.Join(fmt.Errorf("failed to validate all items in array: %v", valItems), err)
		}
	}

	return val, nil
}

func evalItems(s *Schema, valItems []any) ([]any, error) {
	if s == nil {
		return nil, errors.New("\tschema is nil, cannot eval Items")
	}

	if s.Items == nil {
		fmt.Printf("no item schema specified\n")
		return valItems, nil
	}

	validItems := []any{}
	errs := []error{}

	for i := range valItems {
		v, err := s.Items.Eval(valItems[i])
		if err != nil {
			errs = append(errs, err)
		} else {
			validItems = append(validItems, v)
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(fmt.Sprintf("\terror: could not validate all items in array: %v", errs))
	}
	return validItems, nil
}

func GetDataType(v interface{}) string {
	switch v.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case float32, float64:
		return "number"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case []bool, []json.Number, []float32, []float64, []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []string:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

package gojson

import (
	"errors"
	"fmt"
	"reflect"
)

type JsonValueType = string

const (
	STRING JsonValueType = "STRING"
	NUMBER JsonValueType = "NUMBER"
	BOOL   JsonValueType = "BOOLEAN"
	NULL   JsonValueType = "NULL"
	OBJECT JsonValueType = "OBJECT"
	ARRAY  JsonValueType = "ARRAY"
)

type JsonValue struct {
	Value     interface{}
	ValueType JsonValueType
}

type numberConverter = func(i float64) interface{}

var numbers = map[reflect.Kind]numberConverter{
	reflect.Int: func(i float64) interface{} {
		return int(i)
	},
	reflect.Int8: func(i float64) interface{} {
		return int8(i)
	},
	reflect.Int16: func(i float64) interface{} {
		return int16(i)
	},
	reflect.Int32: func(i float64) interface{} {
		return int32(i)
	},
	reflect.Int64: func(i float64) interface{} {
		return int64(i)
	},
	reflect.Uint: func(i float64) interface{} {
		return uint(i)
	},
	reflect.Uint8: func(i float64) interface{} {
		return uint8(i)
	},
	reflect.Uint16: func(i float64) interface{} {
		return uint16(i)
	},
	reflect.Uint32: func(i float64) interface{} {
		return uint32(i)
	},
	reflect.Uint64: func(i float64) interface{} {
		return uint64(i)
	},
	reflect.Float32: func(i float64) interface{} {
		return float32(i)
	},
	reflect.Float64: func(i float64) interface{} {
		return i
	},
}

var supportedKinds = map[reflect.Kind]JsonValueType{
	reflect.Bool:   BOOL,
	reflect.String: STRING,
	reflect.Slice:  ARRAY,
	reflect.Map:    OBJECT,
	reflect.Struct: OBJECT,
}

func isSupported(k reflect.Kind) (JsonValueType, bool) {
	if jt, ok := supportedKinds[k]; ok {
		return jt, true
	}

	if _, ok := numbers[k]; ok {
		return NUMBER, true
	}

	return "", false
}

func (jv *JsonValue) Unmarshal(obj any) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Pointer {
		return errors.New("expected: a pointer")
	}

	kind := v.Elem().Kind()

	if _, ok := isSupported(kind); !ok {
		return errors.New(fmt.Sprintf("unsupported type: %s", kind.String()))
	}

	return jv.setValue(kind, v.Elem())
}

func (jv *JsonValue) setValue(kind reflect.Kind, v reflect.Value) error {
	jt, _ := isSupported(kind)
	if jt != jv.ValueType {
		return errors.New(fmt.Sprintf("type mismatch: expected: %s, provided: %s", jv.ValueType, jt))
	}

	if kind == reflect.String {
		v.Set(reflect.ValueOf(jv.Value))
	} else if kind == reflect.Bool {
		v.Set(reflect.ValueOf(jv.Value))
	} else if converter, ok := numbers[kind]; ok {
		v.Set(reflect.ValueOf(converter(jv.Value.(float64))))
	} else if kind == reflect.Slice {
		dataType := v.Type().Elem().Kind()

		values := jv.Value.([]JsonValue)
		var jsonType = values[0].ValueType

		for _, value := range values {
			if value.ValueType != jsonType {
				return errors.New("json array does not have elements of one type")
			}
		}

		if jt, ok = isSupported(dataType); !ok || jt != jsonType {
			return errors.New("type mismatch for array")
		}

		refSlice := reflect.MakeSlice(reflect.SliceOf(v.Type().Elem()), len(values), len(values))

		for i := 0; i < len(values); i++ {
			if err := values[i].setValue(dataType, refSlice.Index(i)); err != nil {
				return err
			}
		}
		v.Set(refSlice)
	} else if kind == reflect.Struct {
		m := jv.Value.(map[string]JsonValue)

		for k, val := range m {
			f := v.FieldByName(k)
			if err := val.setValue(f.Kind(), f); err != nil {
				return err
			}
		}

	}
	return nil
}

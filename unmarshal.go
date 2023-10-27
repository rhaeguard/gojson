package gojson

type JsonValueType = string

const (
	STRING  JsonValueType = "STRING"
	NUMBER  JsonValueType = "NUMBER"
	BOOL    JsonValueType = "BOOLEAN"
	NULL    JsonValueType = "NULL"
	OBJECT  JsonValueType = "OBJECT"
	ARRAY   JsonValueType = "ARRAY"
	UNKNOWN JsonValueType = "UNKNOWN"
)

type JsonValue struct {
	Value     interface{}
	ValueType JsonValueType
}

func (v JsonValue) Type() JsonValueType {
	return v.ValueType
}

package gojson

import (
	"fmt"
	"strconv"
	"strings"
)

type elementType = string

const (
	number   elementType = "<number>"
	integer  elementType = "<integer>"
	value    elementType = "<value>"
	array    elementType = "<array>"
	members  elementType = "<object fields>"
	member   elementType = "<object field>"
	element  elementType = "<array element>"
	elements elementType = "<array elements>"
	object   elementType = "<object>"
	boolean  elementType = "<boolean>"
	exponent elementType = "<exponent>"
	fraction elementType = "<fraction>"
	/* the rest represents literal tokens */
	ltObjectStart    elementType = "{"
	ltObjectEnd      elementType = "}"
	ltArrayStart     elementType = "["
	ltArrayEnd       elementType = "]"
	ltComma          elementType = ","
	ltColon          elementType = ":"
	ltFractionSymbol elementType = "."
	ltBoolean        elementType = "<bool_literal>"
	ltExponent       elementType = "e/E"
	ltDigits         elementType = "[0-9] (digits)"
	ltNull           elementType = "<null>"
	ltSign           elementType = "+/-"
	ltString         elementType = "<string_literal>"
)

type grammarRule struct {
	lhs    string
	rhs    [][]elementType
	toJson func(values ...*stackElement) JsonValue
}

var grammar = []grammarRule{
	grammarRule{value, [][]elementType{
		{object},
		{array},
		{number},
		{boolean},
		{ltString},
		{ltNull},
	}, func(values ...*stackElement) JsonValue {
		v := values[0].Value()
		if str, ok := v.(string); ok {
			return JsonValue{
				Value:     str,
				ValueType: STRING,
			}
		} else if v == nil {
			return JsonValue{
				Value:     nil,
				ValueType: NULL,
			}
		}
		return values[0].asJsonValue()
	}},
	grammarRule{boolean, [][]elementType{
		{ltBoolean},
	}, func(values ...*stackElement) JsonValue {
		b := values[0].Value().(string)
		return JsonValue{
			Value:     b == "true",
			ValueType: BOOL,
		}
	}},
	grammarRule{object, [][]elementType{
		{ltObjectStart, ltObjectEnd},
		{ltObjectStart, members, ltObjectEnd},
	}, func(values ...*stackElement) JsonValue {
		// TODO: incomplete
		if len(values) == 2 {
			return JsonValue{
				Value:     map[string]JsonValue{},
				ValueType: OBJECT,
			}
		}
		return values[1].asJsonValue()
	}},
	grammarRule{members, [][]elementType{
		{member},
		{members, ltComma, member},
	}, func(values ...*stackElement) JsonValue {
		size := len(values)
		members := map[string]JsonValue{}
		member := values[size-1].asJsonValue().Value.(map[string]JsonValue)

		if size == 3 {
			members = values[0].asJsonValue().Value.(map[string]JsonValue)
		}

		for k, v := range member {
			members[k] = v
		}

		return JsonValue{
			ValueType: OBJECT,
			Value:     members,
		}
	}},
	grammarRule{member, [][]elementType{
		{ltString, ltColon, value},
	}, func(values ...*stackElement) JsonValue {
		key := fmt.Sprintf("%s", values[0].Value())
		valueObj := values[2].asJsonValue()

		return JsonValue{
			ValueType: OBJECT,
			Value: map[string]JsonValue{
				key: valueObj,
			},
		}
	}},
	grammarRule{array, [][]elementType{
		{ltArrayStart, ltArrayEnd},
		{ltArrayStart, elements, ltArrayEnd},
	}, func(values ...*stackElement) JsonValue {
		if len(values) == 2 {
			return JsonValue{
				ValueType: ARRAY,
				Value:     []JsonValue{},
			}
		}
		return values[1].asJsonValue()
	}},
	grammarRule{elements, [][]elementType{
		{element},
		{elements, ltComma, element},
	}, func(values ...*stackElement) JsonValue {
		size := len(values)

		var elements []JsonValue
		if size == 3 {
			elements = (values[0].asJsonValue()).Value.([]JsonValue)
		}

		element := values[size-1].asJsonValue()
		elements = append(elements, element)

		return JsonValue{
			ValueType: ARRAY,
			Value:     elements,
		}
	}},
	grammarRule{element, [][]elementType{
		{value},
	}, func(values ...*stackElement) JsonValue {
		return values[0].asJsonValue()
	}},
	grammarRule{number, [][]elementType{
		{integer, fraction, exponent},
		{integer, fraction},
		{integer, exponent},
		{integer},
	}, func(values ...*stackElement) JsonValue {
		size := len(values)
		var integerValue = values[0].asJsonValue().Value.(string)

		var fraction string
		if size >= 2 && strings.HasPrefix(values[1].asJsonValue().Value.(string), ".") {
			fraction = values[1].asJsonValue().Value.(string)
		} else {
			fraction = ""
		}

		var exponent string
		if size == 2 && strings.HasPrefix(values[1].asJsonValue().Value.(string), "e") {
			exponent = values[1].asJsonValue().Value.(string)
		} else if size == 3 && strings.HasPrefix(values[2].asJsonValue().Value.(string), "e") {
			exponent = values[2].asJsonValue().Value.(string)
		} else {
			exponent = ""
		}

		expression := fmt.Sprintf("%s%s%s", integerValue, fraction, exponent)
		value, err := strconv.ParseFloat(expression, 64) // TODO: potential for an error!

		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}

		return JsonValue{
			Value:     value,
			ValueType: NUMBER,
		}
	}},
	grammarRule{integer, [][]elementType{
		{ltDigits},
		{ltSign, ltDigits},
	}, func(values ...*stackElement) JsonValue {
		size := len(values)
		digits := values[size-1]
		var sign uint8 = '+'
		if size == 2 {
			sign = values[0].Value().(uint8) // - or +
		}
		v := fmt.Sprintf("%c%s", sign, digits.Value())
		return JsonValue{
			Value:     v,
			ValueType: NUMBER,
		}
	}},
	grammarRule{fraction, [][]elementType{
		{ltFractionSymbol, ltDigits},
	}, func(values ...*stackElement) JsonValue {
		var fractionDigits = fmt.Sprintf(".%s", values[1].Value())

		return JsonValue{
			Value:     fractionDigits,
			ValueType: NUMBER,
		}
	}},
	grammarRule{exponent, [][]elementType{
		{ltExponent, integer},
	}, func(values ...*stackElement) JsonValue {
		var exponentExpr = fmt.Sprintf("e%s", values[1].asJsonValue().Value.(string))

		return JsonValue{
			Value:     exponentExpr,
			ValueType: NUMBER,
		}
	}},
}

type jsonElement struct {
	value           interface{}
	jsonElementType elementType
}

type stackElement struct {
	value token
	rule  *jsonElement
}

func (se stackElement) String() string {
	if se.rule == nil {
		return fmt.Sprintf("%s", se.value.tokenType)
	}
	return fmt.Sprintf("%s", se.rule.jsonElementType)
}

func (se stackElement) Value() interface{} {
	if se.rule == nil {
		return se.value.value
	}
	return se.rule.value
}

func (se stackElement) asJsonValue() JsonValue {
	return se.rule.value.(JsonValue)
}

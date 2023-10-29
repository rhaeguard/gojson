package gojson

import (
	"fmt"
	"strconv"
	"strings"
)

type ElementType = string

const (
	Number   ElementType = "NUMBER"
	Integer  ElementType = "INTEGER"
	Value    ElementType = "VALUE"
	Array    ElementType = "ARRAY"
	Members  ElementType = "MEMBERS"
	Member   ElementType = "MEMBER"
	Element  ElementType = "ELEMENT"
	Elements ElementType = "ELEMENTS"
	Object   ElementType = "OBJECT"
	Boolean  ElementType = "BOOLEAN"
	Exponent ElementType = "EXPONENT"
	Fraction ElementType = "FRACTION"
	/* the rest represents literal tokens */
	TTObjectStart    ElementType = "TT_OBJECT_START"
	TTObjectEnd      ElementType = "TT_OBJECT_END"
	TTArrayStart     ElementType = "TT_ARRAY_START"
	TTArrayEnd       ElementType = "TT_ARRAY_END"
	TTComma          ElementType = "TT_COMMA"
	TTColon          ElementType = "TT_COLON"
	TTFractionSymbol ElementType = "TT_FRACTION_SYMBOL"
	TTBoolean        ElementType = "TT_BOOLEAN"
	TTExponent       ElementType = "TT_EXPONENT"
	TTDigits         ElementType = "TT_DIGITS"
	TTNull           ElementType = "TT_NULL"
	TTSign           ElementType = "TT_SIGN"
	TTString         ElementType = "TT_STRING"
)

type GrammarRule struct {
	lhs    string
	rhs    [][]ElementType
	toJson func(values ...*StackElement) JsonValue
}

var newGrammar = []GrammarRule{
	GrammarRule{Value, [][]ElementType{
		{Object},
		{Array},
		{Number},
		{Boolean},
		{TTString},
		{TTNull},
	}, func(values ...*StackElement) JsonValue {
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
		return values[0].Value().(JsonValue)
	}},
	GrammarRule{Boolean, [][]ElementType{
		{TTBoolean},
	}, func(values ...*StackElement) JsonValue {
		b := values[0].Value().(string)
		return JsonValue{
			Value:     b == "true",
			ValueType: BOOL,
		}
	}},
	GrammarRule{Object, [][]ElementType{
		{TTObjectStart, TTObjectEnd},
		{TTObjectStart, Members, TTObjectEnd},
	}, func(values ...*StackElement) JsonValue {
		// TODO: incomplete
		if len(values) == 2 {
			return JsonValue{
				Value:     map[string]JsonValue{},
				ValueType: OBJECT,
			}
		}
		return values[1].Value().(JsonValue)
	}},
	GrammarRule{Members, [][]ElementType{
		{Member},
		{Members, TTComma, Member},
	}, func(values ...*StackElement) JsonValue {
		size := len(values)
		var members = map[string]JsonValue{}
		member := values[size-1].Value().(JsonValue).Value.(map[string]JsonValue)

		if size == 3 {
			members = values[0].Value().(JsonValue).Value.(map[string]JsonValue)
		}

		for k, v := range member {
			members[k] = v
		}

		return JsonValue{
			ValueType: OBJECT,
			Value:     members,
		}
	}},
	GrammarRule{Member, [][]ElementType{
		{TTString, TTColon, Value},
	}, func(values ...*StackElement) JsonValue {
		keyName := values[0]
		valueObj := values[2].Value().(JsonValue)

		key := fmt.Sprintf("%s", keyName.Value())

		return JsonValue{
			ValueType: OBJECT,
			Value: map[string]JsonValue{
				key: valueObj,
			},
		}
	}},
	GrammarRule{Array, [][]ElementType{
		{TTArrayStart, TTArrayEnd},
		{TTArrayStart, Elements, TTArrayEnd},
	}, func(values ...*StackElement) JsonValue {
		if len(values) == 2 {
			return JsonValue{
				ValueType: ARRAY,
				Value:     []JsonValue{},
			}
		}
		return values[1].Value().(JsonValue)
	}},
	GrammarRule{Elements, [][]ElementType{
		{Element},
		{Elements, TTComma, Element},
	}, func(values ...*StackElement) JsonValue {
		size := len(values)

		var elements []JsonValue
		if size == 3 {
			elements = (values[0].Value().(JsonValue)).Value.([]JsonValue)
		}

		element := values[size-1].Value().(JsonValue)
		elements = append(elements, element)

		return JsonValue{
			ValueType: ARRAY,
			Value:     elements,
		}
	}},
	GrammarRule{Element, [][]ElementType{
		{Value},
	}, func(values ...*StackElement) JsonValue {
		return values[0].Value().(JsonValue)
	}},
	GrammarRule{Number, [][]ElementType{
		{Integer, Fraction, Exponent},
		{Integer, Fraction},
		{Integer, Exponent},
		{Integer},
	}, func(values ...*StackElement) JsonValue {
		size := len(values)
		var integerValue = values[0].Value().(JsonValue).Value.(string)

		var fraction string
		if size >= 2 && strings.HasPrefix(values[1].Value().(JsonValue).Value.(string), ".") {
			fraction = values[1].Value().(JsonValue).Value.(string)
		} else {
			fraction = ""
		}

		var exponent string
		if size == 2 && strings.HasPrefix(values[1].Value().(JsonValue).Value.(string), "e") {
			exponent = values[1].Value().(JsonValue).Value.(string)
		} else if size == 3 && strings.HasPrefix(values[2].Value().(JsonValue).Value.(string), "e") {
			exponent = values[2].Value().(JsonValue).Value.(string)
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
	GrammarRule{Integer, [][]ElementType{
		{TTDigits},
		{TTSign, TTDigits},
	}, func(values ...*StackElement) JsonValue {
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
	GrammarRule{Fraction, [][]ElementType{
		{TTFractionSymbol, TTDigits},
	}, func(values ...*StackElement) JsonValue {
		var fractionDigits = fmt.Sprintf(".%s", values[1].Value())

		return JsonValue{
			Value:     fractionDigits,
			ValueType: NUMBER,
		}
	}},
	GrammarRule{Exponent, [][]ElementType{
		{TTExponent, Integer},
	}, func(values ...*StackElement) JsonValue {
		var exponentExpr = fmt.Sprintf("e%s", values[1].Value().(JsonValue).Value.(string))

		return JsonValue{
			Value:     exponentExpr,
			ValueType: NUMBER,
		}
	}},
}

type JsonElement struct {
	value           interface{}
	jsonElementType ElementType
}

type StackElement struct {
	value Token
	rule  *JsonElement
}

func (se StackElement) String() string {
	if se.rule == nil {
		return fmt.Sprintf("%s", se.value.tokenType)
	}
	return fmt.Sprintf("%s", se.rule.jsonElementType)
}

func (se StackElement) Value() interface{} {
	if se.rule == nil {
		return se.value.value
	}
	return se.rule.value
}

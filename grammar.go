package gojson

import (
	"fmt"
	"sort"
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
	Lhs    string
	Rhs    [][]ElementType
	ToJson func(values ...*StackElement) JsonValue
}

func grammarRule(
	lhs string,
	rhs [][]ElementType,
	toJson func(values ...*StackElement) JsonValue,
) GrammarRule {
	return GrammarRule{
		Lhs:    lhs,
		Rhs:    rhs,
		ToJson: toJson,
	}
}

var newGrammar = []GrammarRule{
	grammarRule(Value, [][]ElementType{
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
	}),
	grammarRule(Boolean, [][]ElementType{
		{TTBoolean},
	}, func(values ...*StackElement) JsonValue {
		b := values[0].Value().(string)
		if b == "true" {
			return JsonValue{
				Value:     true,
				ValueType: BOOL,
			}
		} else if b == "false" {
			return JsonValue{
				Value:     false,
				ValueType: BOOL,
			}
		}

		return JsonValue{
			Value:     nil,
			ValueType: UNKNOWN,
		}
	}),
	grammarRule(Object, [][]ElementType{
		{TTObjectStart, TTObjectEnd},
		{TTObjectStart, Members, TTObjectEnd},
	}, func(values ...*StackElement) JsonValue {
		// TODO: incomplete
		if len(values) == 3 {
			return values[1].Value().(JsonValue)
		}
		return JsonValue{
			Value: map[string]JsonValue{},
		}
	}),
	grammarRule(Members, [][]ElementType{
		{Member},
		{Members, TTComma, Member},
	}, func(values ...*StackElement) JsonValue {
		if len(values) == 1 {
			return values[0].Value().(JsonValue)
		} else if len(values) == 3 {
			mp := values[0].Value().(JsonValue)
			n := values[2].Value().(JsonValue)

			kvs := mp.Value.(map[string]JsonValue)
			kvn := n.Value.(map[string]JsonValue)

			for k, v := range kvn {
				kvs[k] = v
			}

			return JsonValue{
				ValueType: OBJECT,
				Value:     kvs,
			}
		}

		return JsonValue{
			ValueType: UNKNOWN,
		}
	}),
	grammarRule(Member, [][]ElementType{
		{TTString, TTColon, Value},
	}, func(values ...*StackElement) JsonValue {
		keyName := values[0]
		valueObj := values[2]

		key := fmt.Sprintf("%s", keyName.Value())

		return JsonValue{
			ValueType: OBJECT,
			Value: map[string]JsonValue{
				key: valueObj.Value().(JsonValue),
			},
		}
	}),
	grammarRule(Array, [][]ElementType{
		{TTArrayStart, TTArrayEnd},
		{TTArrayStart, Elements, TTArrayEnd},
	}, func(values ...*StackElement) JsonValue {
		if len(values) == 2 {
			return JsonValue{
				ValueType: ARRAY,
				Value:     []JsonValue{},
			}
		} else if len(values) == 3 {
			return values[1].Value().(JsonValue)
		}
		return JsonValue{
			ValueType: UNKNOWN,
		}
	}),
	grammarRule(Elements, [][]ElementType{
		{Element},
		{Elements, TTComma, Element},
	}, func(values ...*StackElement) JsonValue {
		// TODO: incomplete
		if len(values) == 1 {
			return JsonValue{
				ValueType: ARRAY,
				Value: []JsonValue{
					values[0].Value().(JsonValue),
				},
			}
		} else if len(values) == 3 {
			jsonValues := (values[0].Value().(JsonValue)).Value.([]JsonValue)
			e := values[2].Value().(JsonValue)
			jsonValues = append(jsonValues, e)

			return JsonValue{
				ValueType: ARRAY,
				Value:     jsonValues,
			}
		}
		return JsonValue{
			ValueType: UNKNOWN,
		}
	}),
	grammarRule(Element, [][]ElementType{
		{Value},
	}, func(values ...*StackElement) JsonValue {
		return values[0].Value().(JsonValue)
	}),
	grammarRule(Number, [][]ElementType{
		{Integer, Fraction, Exponent},
		{Integer, Fraction},
		{Integer, Exponent},
		{Integer},
	}, func(values ...*StackElement) JsonValue {
		size := len(values)
		var integerValue = values[0].Value().(JsonValue).Value.(string)

		var fractionDigits string
		if size == 2 && strings.HasPrefix(values[1].Value().(JsonValue).Value.(string), ".") {
			fractionDigits = values[1].Value().(JsonValue).Value.(string)
		} else {
			fractionDigits = ""
		}

		var exponentDigits string
		if size == 2 && strings.HasPrefix(values[1].Value().(JsonValue).Value.(string), "e") {
			exponentDigits = values[1].Value().(JsonValue).Value.(string)
		} else if (size == 3 && fractionDigits != "") && strings.HasPrefix(values[2].Value().(JsonValue).Value.(string), "e") {
			exponentDigits = values[2].Value().(JsonValue).Value.(string)
		} else {
			exponentDigits = ""
		}

		expression := fmt.Sprintf("%s%s%s", integerValue, fractionDigits, exponentDigits)
		value, _ := strconv.ParseFloat(expression, 64)

		return JsonValue{
			Value:     value,
			ValueType: NUMBER,
		}
	}),
	grammarRule(Integer, [][]ElementType{
		{TTDigits},
		{TTSign, TTDigits},
	}, func(values ...*StackElement) JsonValue {
		size := len(values)
		if size == 1 {
			// TODO: handle errors properly
			v := fmt.Sprintf("%s", values[0].Value())
			return JsonValue{
				Value:     v,
				ValueType: NUMBER,
			}
		} else if size == 2 {
			signStr := values[0]
			digits := values[1]
			v := fmt.Sprintf("%c%s", signStr.Value().(uint8), digits.Value())
			return JsonValue{
				Value:     v,
				ValueType: NUMBER,
			}
		}
		return JsonValue{
			ValueType: UNKNOWN,
		}
	}),
	grammarRule(Fraction, [][]ElementType{
		{TTFractionSymbol, TTDigits},
	}, func(values ...*StackElement) JsonValue {
		if len(values) != 2 {
			return JsonValue{
				ValueType: UNKNOWN,
			}
		}

		var fractionDigits = fmt.Sprintf(".%s", values[1].Value())

		return JsonValue{
			Value:     fractionDigits,
			ValueType: NUMBER,
		}
	}),
	grammarRule(Exponent, [][]ElementType{
		{TTExponent, Integer},
	}, func(values ...*StackElement) JsonValue {
		if len(values) != 2 {
			return JsonValue{
				ValueType: UNKNOWN,
			}
		}

		var exponentExpr = fmt.Sprintf("e%s", values[1].Value())

		return JsonValue{
			Value:     exponentExpr,
			ValueType: NUMBER,
		}
	}),
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

func anyIncompletePrefix(candidates ...ElementType) (string, bool) {
	// find all matches
	// full or partial
	// only match or multiple matches
	type payload struct {
		matchType string
		prodSize  int
	}
	data := []payload{}
	for _, rule := range newGrammar {
		outcomes := rule.Rhs
		for _, production := range outcomes {
			cSize := len(candidates)
			rSize := len(production)

			if cSize > rSize {
				continue
			}

			continueOuter := false
			for i := 0; i < cSize; i++ {
				if candidates[i] != production[i] {
					continueOuter = true
				}
			}
			if continueOuter {
				continue
			}

			var p payload
			if cSize == rSize {
				p = payload{
					matchType: "full",
					prodSize:  rSize,
				}
			} else {
				p = payload{
					matchType: "partial",
					prodSize:  rSize,
				}
			}
			data = append(data, p)
		}
	}

	if len(data) == 0 {
		return "none", false
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].prodSize > data[j].prodSize
	})

	return data[0].matchType, true
}

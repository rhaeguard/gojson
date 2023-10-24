package main

import (
	"fmt"
	"strconv"
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
	ToJson func(values ...*StackElement) interface{}
}

func grammarRule(lhs string, rhs [][]ElementType) GrammarRule {
	return GrammarRule{
		Lhs: lhs,
		Rhs: rhs,
		ToJson: func(values ...*StackElement) interface{} {
			return nil
		},
	}
}

func grammarRule2(
	lhs string,
	rhs [][]ElementType,
	toJson func(values ...*StackElement) interface{},
) GrammarRule {
	return GrammarRule{
		Lhs:    lhs,
		Rhs:    rhs,
		ToJson: toJson,
	}
}

var newGrammar = []GrammarRule{
	grammarRule2(Value, [][]ElementType{
		{Object},
		{Array},
		{Number},
		{TTString},
		{TTBoolean},
		{TTNull},
	}, func(values ...*StackElement) interface{} {
		return values[0].Value()
	}),
	grammarRule2(Object, [][]ElementType{
		{TTObjectStart, TTObjectEnd},
		{TTObjectStart, Members, TTObjectEnd},
	}, func(values ...*StackElement) interface{} {
		// TODO: incomplete
		if len(values) == 3 {
			return values[1].Value()
		}
		return nil
	}),
	grammarRule2(Members, [][]ElementType{
		{Member},
		{Members, TTComma, Member},
	}, func(values ...*StackElement) interface{} {
		if len(values) == 1 {
			return values[0].Value()
		} else if len(values) == 3 {
			mp := values[0].Value().(map[string]interface{})
			n := values[2].Value().(map[string]interface{})
			for k, v := range n {
				mp[k] = v
			}
			return mp
		}

		return nil
	}),
	grammarRule2(Member, [][]ElementType{
		{TTString, TTColon, Value},
	}, func(values ...*StackElement) interface{} {
		keyName := values[0]
		valueObj := values[2]

		key := fmt.Sprintf("%s", keyName.Value())

		return map[string]interface{}{
			key: valueObj.Value(),
		}
	}),
	grammarRule(Array, [][]ElementType{
		{TTArrayStart, TTArrayEnd},
		{TTArrayStart, Elements, TTArrayEnd},
	}),
	grammarRule2(Elements, [][]ElementType{
		{Element},
		{Elements, TTComma, Element},
	}, func(values ...*StackElement) interface{} {
		// TODO: incomplete
		return values[0].Value()
	}),
	grammarRule2(Element, [][]ElementType{
		{Value},
	}, func(values ...*StackElement) interface{} {
		return values[0].Value()
	}),
	grammarRule2(Number, [][]ElementType{
		{Integer, TTFractionSymbol, TTDigits, TTExponent, Integer},
		{Integer, TTFractionSymbol, TTDigits},
		{Integer, TTExponent, Integer},
		{Integer},
	}, func(values ...*StackElement) interface{} {
		// TODO: handle all cases
		if len(values) >= 1 {
			integer := values[0]
			return integer.Value()
		}
		return nil
	}),
	grammarRule2(Integer, [][]ElementType{
		{TTDigits},
		{TTSign, TTDigits},
	}, func(values ...*StackElement) interface{} {
		s := len(values)
		if s == 1 {
			digits := values[0]
			// TODO: handle errors properly
			v, _ := strconv.Atoi(fmt.Sprintf("%s", digits.Value()))
			return v
		} else if s == 2 {
			signStr := values[0]
			digits := values[1]
			// TODO: handle errors properly
			sign, _ := strconv.Atoi(fmt.Sprintf("%s", signStr.Value()))
			v, _ := strconv.Atoi(fmt.Sprintf("%s", digits.Value()))
			return sign * v
		}
		return nil
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

func anyPartialMatch(candidates ...ElementType) (ElementType, bool) {
	// find all matches
	// if the only match is a full-match => return appropriately
	// otherwise it's a partial match => return appropriately
	// if no match => false
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

			if cSize != rSize {
			}
			return rule.Lhs, true
		}
	}

	return "", false
}

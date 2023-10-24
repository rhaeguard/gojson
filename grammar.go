package main

import (
	"fmt"
)

type JsonElementType = string

const (
	Number     JsonElementType = "NUMBER"
	Integer    JsonElementType = "INTEGER"
	Exponent   JsonElementType = "EXPONENT"
	Digits     JsonElementType = "DIGITS"
	Digit      JsonElementType = "DIGIT"
	OneNine    JsonElementType = "ONE_NINE"
	Sign       JsonElementType = "SIGN"
	String     JsonElementType = "STRING"
	JsonValue  JsonElementType = "JSON"
	Element    JsonElementType = "ELEMENT"
	Value      JsonElementType = "VALUE"
	Array      JsonElementType = "ARRAY"
	Whitespace JsonElementType = "WHITESPACE"
	Members    JsonElementType = "MEMBERS"
	Member     JsonElementType = "MEMBER"
	Elements   JsonElementType = "ELEMENTS"
	Character  JsonElementType = "CHARACTER"
	Characters JsonElementType = "CHARACTERS"
	Object     JsonElementType = "OBJECT"
)

func generateChars(start, end uint8) [][]JsonElementType {
	var result [][]JsonElementType
	for ch := start; ch <= end; ch++ {
		a := []JsonElementType{
			fmt.Sprintf("%c", ch),
		}
		result = append(result, a)
	}

	return result
}

type GrammarRule struct {
	Lhs string
	Rhs [][]JsonElementType
}

func grammarRule(lhs string, rhs [][]JsonElementType) GrammarRule {
	return GrammarRule{
		Lhs: lhs,
		Rhs: rhs,
	}
}

var newGrammar = []GrammarRule{
	grammarRule(JsonValue, [][]JsonElementType{
		{Element},
	}),
	grammarRule(Value, [][]JsonElementType{
		{Object},
		{Array},
		{Number},
		{TTString},
		{TTBoolean},
		{TTNull},
	}),
	grammarRule(Object, [][]JsonElementType{
		{TTObjectStart, TTObjectEnd},
		{TTObjectStart, Members, TTObjectEnd},
	}),
	grammarRule(Members, [][]JsonElementType{
		{Member},
		{Members, TTComma, Member},
	}),
	grammarRule(Member, [][]JsonElementType{
		{TTString, TTColon, Value},
	}),
	grammarRule(Array, [][]JsonElementType{
		{TTArrayStart, TTArrayEnd},
		{TTArrayStart, Elements, TTArrayEnd},
	}),
	grammarRule(Elements, [][]JsonElementType{
		{Element},
		{Elements, TTComma, Element},
	}),
	grammarRule(Element, [][]JsonElementType{
		{Value},
	}),
	grammarRule(Number, [][]JsonElementType{
		{Integer, TTFractionSymbol, TTDigits, TTExponent, Integer},
		{Integer, TTFractionSymbol, TTDigits},
		{Integer, TTExponent, Integer},
		{Integer},
	}),
	grammarRule(Integer, [][]JsonElementType{
		{TTDigits},
		{TTSign, TTDigits},
	}),
}

type JsonElement struct {
	value           interface{}
	jsonElementType JsonElementType
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

func anyPartialMatch(candidates ...JsonElementType) (JsonElementType, bool) {
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
				return rule.Lhs, true
			}
		}
	}

	return "", false
}

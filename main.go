package main

import (
	"fmt"
	"strconv"
	"strings"
)

const ExampleJson = `{"value":1,"name":"renault"}`

type JsonElementType = string

const (
	Number   JsonElementType = "NUMBER"
	String   JsonElementType = "STRING"
	Object   JsonElementType = "OBJECT"
	KeyValue JsonElementType = "KEY_VALUE"
)

type JsonElement struct {
	value           interface{}
	jsonElementType JsonElementType
}

type StackElement struct {
	value string
	rule  *JsonElement
}

func (se *StackElement) is(value string) bool {
	return value == se.value
}

func (se *StackElement) isOfType(value JsonElementType) bool {
	return se.rule != nil && value == se.rule.jsonElementType
}

func (se *StackElement) isOfEitherType(values ...JsonElementType) bool {
	for _, v := range values {
		if se.isOfType(v) {
			return true
		}
	}
	return false
}

// returns the element and the number of stack elements
// to be removed from the top of the stack
func numericValue(stack []*StackElement) (*JsonElement, int) {
	if len(stack) < 1 {
		return nil, 0
	}
	el := stack[len(stack)-1]
	if el.rule == nil {
		if intValue, err := strconv.Atoi(el.value); err == nil {
			return &JsonElement{
				value:           intValue,
				jsonElementType: Number,
			}, 1
		}
	}
	return nil, 0
}

func keyValuePair(stack []*StackElement) (*JsonElement, int) {
	// string : value
	if len(stack) < 3 {
		return nil, 0
	}

	elements := stack[len(stack)-3:]

	if !elements[1].is(`:`) {
		return nil, 0
	}

	if !elements[0].isOfType(String) {
		return nil, 0
	}

	if !elements[2].isOfEitherType(Number, String, Object) {
		return nil, 0
	}

	key := elements[0].rule.value.(string)
	value := elements[2].rule.value

	return &JsonElement{
		value: map[string]interface{}{
			key: value,
		},
		jsonElementType: KeyValue,
	}, 3

}

func objectValue(stack []*StackElement) (*JsonElement, int) {
	// {"value":1}
	// { keyValuePair (, keyValuePairs) }
	if len(stack) < 1 {
		return nil, 0
	}
	maybeClosingBrace := stack[len(stack)-1]
	if !maybeClosingBrace.is(`}`) {
		return nil, 0
	}

	i := len(stack) - 2
	for i >= 0 {
		el := stack[i]
		if el.is(`{`) {
			break
		}
		i--
	}

	if i < 0 || !stack[i].is(`{`) {
		return nil, 0
	}
	// xxxxx{a,b,c,d}
	// 14
	// 5
	// 14 - 5 = 9
	elements := stack[i+1 : len(stack)-1]
	if !elements[0].isOfType(KeyValue) {
		return nil, 0
	}
	var keyValuePairs = []*JsonElement{
		elements[0].rule,
	}

	if len(elements) > 1 {
		for t := 1; t < len(elements)-1; t++ {
			comma := elements[t]
			kvpair := elements[t+1]

			if !comma.is(`,`) {
				return nil, 0
			}

			if !kvpair.isOfType(KeyValue) {
				return nil, 0
			}

			keyValuePairs = append(keyValuePairs, kvpair.rule)
		}
	}

	return &JsonElement{
		value:           keyValuePairs,
		jsonElementType: Object,
	}, len(elements) + 2
}

func stringValue(stack []*StackElement) (*JsonElement, int) {
	if len(stack) < 2 { // at least 2 because empty string: ""
		return nil, 0
	}
	last := stack[len(stack)-1]
	if !last.is(`"`) {
		return nil, 0
	}

	i := len(stack) - 2
	for i >= 0 {
		el := stack[i]
		if el.is(`"`) {
			break
		}
		i--
	}

	if i < 0 || !stack[i].is(`"`) {
		return nil, 0
	}

	var sb strings.Builder
	elements := stack[i+1 : len(stack)-1]
	for _, se := range elements {
		sb.WriteString(se.value)
	}
	// abc"value"
	// 10
	// 10 - 4
	return &JsonElement{
		value:           sb.String(),
		jsonElementType: String,
	}, len(stack) - i
}

type Action = func(stack []*StackElement) (*JsonElement, int)

var actions = []Action{
	stringValue,
	numericValue,
	keyValuePair,
	objectValue,
}

func main() {
	var stack []*StackElement

	for _, ch := range ExampleJson {
		stack = append(stack, &StackElement{
			value: fmt.Sprintf("%c", ch),
		})
		continueReduction := true
		for continueReduction {
			continueReduction = false
			for _, action := range actions {
				jsonElement, offset := action(stack)
				if offset != 0 {
					stack = stack[:len(stack)-offset]
					stack = append(stack, &StackElement{
						rule: jsonElement,
					})
					fmt.Printf("%v\n", jsonElement)
					continueReduction = true
					break
				}
			}
		}
	}

	fmt.Printf("%-v\n", stack)
}

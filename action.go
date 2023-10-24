package main

import "fmt"

func action(stack []*StackElement) (*JsonElement, int) {
	stackSize := len(stack)

	var je *JsonElement
	var offsetSize int

	for _, rule := range newGrammar {
		lhs := rule.Lhs
		expansions := rule.Rhs
		for _, expansion := range expansions {
			size := len(expansion)
			if size > stackSize {
				continue
			}
			actual, values := topNOfStack(stack, size)
			matches := compare(expansion, actual)
			if matches && size > offsetSize { // TODO: what if they are equal?
				je = &JsonElement{
					value:           values,
					jsonElementType: lhs,
				}
				offsetSize = size
			}
		}
	}
	return je, offsetSize
}

func topNOfStack(stack []*StackElement, count int) ([]JsonElementType, []string) {
	slice := stack[len(stack)-count:]
	var elements []JsonElementType
	var elementValue []string

	for _, el := range slice {
		var value = el.value.tokenType
		var literal = el.value.value
		if el.rule != nil {
			value = el.rule.jsonElementType
			literal = fmt.Sprintf("%s", el.rule.value)
		}
		elements = append(elements, value)
		elementValue = append(elementValue, fmt.Sprintf("%s", literal))
	}

	return elements, elementValue
}

func compare(expansion, actual []JsonElementType) bool {
	for i := 0; i < len(expansion); i++ {
		if expansion[i] != actual[i] {
			return false
		}
	}
	return true
}

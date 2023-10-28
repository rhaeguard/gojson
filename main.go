package gojson

import "fmt"

func stackToToken(stack []*StackElement) []ElementType {
	var a []ElementType
	for _, e := range stack {
		if e.rule == nil {
			a = append(a, e.value.tokenType)
		} else {
			a = append(a, e.rule.jsonElementType)
		}
	}
	return a
}

func prefixCheck(stack []*StackElement, lookahead Token) string {
	var checkedElements []ElementType

	stackSize := len(stack)
	if stackSize >= 2 {
		checkedElements = append(checkedElements, stackToToken(stack[stackSize-2:])...)
	} else if stackSize == 1 {
		checkedElements = append(checkedElements, stackToToken(stack[0:1])...)
	}

	checkedElements = append(checkedElements, lookahead.tokenType)

	size := len(checkedElements)
	for i := size - 1; i >= 0; i-- {
		if matchType, matches := anyIncompletePrefix(checkedElements[i:size]...); matches {
			return matchType
		}
	}

	return ""
}

func parseJson(input string) (JsonValue, *SyntaxError) {
	tokens, err := lex(input)

	if err != nil {
		return JsonValue{}, err
	}

	var stack []*StackElement

	size := len(tokens)
	noReducePreviously := false

	for i := 0; i < size; {
		lookahead := tokens[i]

		if matchType := prefixCheck(stack, lookahead); matchType != "" {
			i++
			stack = append(stack, &StackElement{value: lookahead})

			if matchType == "partial" {
				continue
			}
			// full match means that there's something we can reduce now
		}

		if noReducePreviously {
			return JsonValue{}, newSyntaxError(-1, fmt.Sprintf("Expected: %s", lookahead.tokenType))
		}

		if jsonElement, offset := action(stack); offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
			noReducePreviously = false
		} else {
			noReducePreviously = true
		}
	}

	for {
		if jsonElement, offset := action(stack); offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
		} else {
			break
		}
	}

	if len(stack) != 1 {
		return JsonValue{}, newSyntaxError(-1, "Parsing failed...")
	}

	values := stack[0].rule.value.(JsonValue).Value.([]JsonValue)

	if len(values) != 1 {
		return JsonValue{}, newSyntaxError(-1, "Parsing failed...")
	}
	return values[0], nil
}

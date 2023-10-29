package gojson

import "fmt"

func ParseJson(input string) (JsonValue, *Error) {
	tokens, err := lex(input)

	if err != nil {
		return JsonValue{}, err
	}

	var stack []*StackElement

	size := len(tokens)
	reducePerformed := true

	for i := 0; i < size; {
		lookahead := tokens[i]

		if matchType := checkIfAnyPrefixExists(stack, lookahead); matchType != noMatch {
			i++
			stack = append(stack, &StackElement{value: lookahead})

			if matchType == partialMatch {
				continue
			}
			// full match means that there's something we can reduce now
		}

		if !reducePerformed {
			return JsonValue{}, newError(-1, fmt.Sprintf("unexpected token: %s", lookahead.tokenType))
		}

		if jsonElement, offset := action(stack); offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
			reducePerformed = true
		} else {
			reducePerformed = false
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
		return JsonValue{}, newError(-1, "parsing failed...")
	}

	values := stack[0].rule.value.(JsonValue).Value.([]JsonValue)
	if len(values) != 1 {
		return JsonValue{}, newError(-1, "parsing failed...")
	}

	return values[0], nil
}

package gojson

func parseJson(input string) (JsonValue, *SyntaxError) {
	tokens, err := lex(input)

	if err != nil {
		return JsonValue{}, err
	}

	var stack []*StackElement

	size := len(tokens)
	for i := 0; i < size; {
		lookahead := tokens[i]

		newStackElement := &StackElement{
			value: lookahead,
		}

		if matchType, matches := anyPartialMatch(lookahead.tokenType); matches {
			i++
			stack = append(stack, newStackElement)

			if matchType == "partial" {
				continue
			}
		} else {
			var topOfStack = ""
			if len(stack) > 1 {
				v := stack[len(stack)-1]
				if v.rule == nil {
					topOfStack = v.value.tokenType
				} else {
					topOfStack = v.rule.jsonElementType
				}
			}

			if topOfStack != "" {
				if matchType, matches := anyPartialMatch(topOfStack, lookahead.tokenType); matches {
					i++
					stack = append(stack, newStackElement)
					if matchType == "partial" {
						continue
					}
				}
			}
		}

		if jsonElement, offset := action(stack); offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
		}
	}

	for {
		continueReduction := false
		jsonElement, offset := action(stack)
		if offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
			continueReduction = true
		}
		if !continueReduction {
			break
		}
	}

	arrayWrapper := stack[0].rule.value.(JsonValue)
	element := arrayWrapper.Value.([]JsonValue)[0]
	return element, nil
}

package main

import (
	"fmt"
)

//const ExampleJson = `{      "value" : 112312312 , "name" : "renault" }`

const ExampleJson = `{ "value" : [1239, 12345], "name" : "renault", "token": true, "hello": null }`

//
//const ExampleJson = `[{ "value" : 12 }]`

func main() {
	tokens := lex(ExampleJson)
	//fmt.Printf("%-v", tokens)
	var stack []*StackElement

	size := len(tokens)
	for i := 0; i < size; {
		lookahead := tokens[i]

		newStackElement := &StackElement{
			value: lookahead,
		}

		fmt.Printf("s: %-v\n", stack)
		if _, matches := anyPartialMatch(lookahead.tokenType); matches {
			i++
			stack = append(stack, newStackElement)
			continue
		}

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
			if _, matches := anyPartialMatch(topOfStack, lookahead.tokenType); matches {
				i++
				stack = append(stack, newStackElement)
				continue
			}
		}

		if jsonElement, offset := action(stack); offset != 0 {
			stack = stack[:len(stack)-offset]
			stack = append(stack, &StackElement{
				rule: jsonElement,
			})
			fmt.Printf("l: %-v\n", stack)
		} else {
			i++
			stack = append(stack, newStackElement)
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
			fmt.Printf("f: %-v\n", stack)
			continueReduction = true
		}
		if !continueReduction {
			break
		}
	}

	fmt.Printf("%-v\n", stack)
}

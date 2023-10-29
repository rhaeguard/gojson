package gojson

import "sort"

type prefixMatch = uint8

const (
	noMatch      prefixMatch = 0
	partialMatch prefixMatch = 1
	fullMatch    prefixMatch = 2
)

// stackToToken - takes a slice of StackElement, and returns a slice of ElementType
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

// checkIfAnyPrefixExists - checks if the combination of max top 2 stack elements
// and the lookahead is a prefix of any known rule in the grammar
func checkIfAnyPrefixExists(stack []*StackElement, lookahead Token) prefixMatch {
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
		if matchType, matches := checkPrefix(checkedElements[i:size]...); matches {
			return matchType
		}
	}

	return noMatch
}

type payload struct {
	matchType uint8
	prodSize  int
}

func checkPrefix(candidates ...ElementType) (prefixMatch, bool) {

	// find all matches
	// full or partial
	// only match or multiple matches
	data := []payload{}
	for _, rule := range newGrammar {
		outcomes := rule.rhs
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
					matchType: fullMatch,
					prodSize:  rSize,
				}
			} else {
				p = payload{
					matchType: partialMatch,
					prodSize:  rSize,
				}
			}
			data = append(data, p)
		}
	}

	if len(data) == 0 {
		return noMatch, false
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].prodSize > data[j].prodSize
	})

	return data[0].matchType, true
}

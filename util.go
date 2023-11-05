package gojson

import "sort"

type prefixMatch = uint8

const (
	noMatch      prefixMatch = 0
	partialMatch prefixMatch = 1
	fullMatch    prefixMatch = 2
)

// stackToToken - takes a slice of stackElement, and returns a slice of ElementType
func stackToToken(stack []*stackElement) []elementType {
	var a []elementType
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
func checkIfAnyPrefixExists(stack []*stackElement, lookahead token) prefixMatch {
	var elems []elementType

	stackSize := len(stack)
	if stackSize >= 2 {
		elems = append(elems, stackToToken(stack[stackSize-2:])...)
	} else if stackSize == 1 {
		elems = append(elems, stackToToken(stack[0:1])...)
	}

	elems = append(elems, lookahead.tokenType)

	size := len(elems)
	for i := size - 1; i >= 0; i-- {
		if matchType := checkPrefix(elems[i:size]...); matchType != noMatch {
			return matchType
		}
	}

	return noMatch
}

type payload struct {
	matchType uint8
	prodSize  int
}

func checkPrefix(candidates ...elementType) prefixMatch {
	// find all matches
	// full or partial
	// only match or multiple matches
	data := []payload{}
	for _, rule := range grammar {
		for _, production := range rule.rhs {
			cSize := len(candidates)
			rSize := len(production)

			// if candidates size is longer than the rule
			// then the rule is not relevant, move on
			if cSize > rSize {
				continue
			}

			didNotMatch := false
			for i := 0; i < cSize; i++ {
				if candidates[i] != production[i] {
					didNotMatch = true
					break
				}
			}
			if didNotMatch {
				continue
			}

			var p payload
			if cSize == rSize {
				p = payload{fullMatch, rSize}
			} else {
				p = payload{partialMatch, rSize}
			}
			data = append(data, p)
		}
	}

	if len(data) == 0 {
		return noMatch
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].prodSize > data[j].prodSize
	})

	return data[0].matchType
}

package gojson

func action(stack []*StackElement) (*JsonElement, int) {
	stackSize := len(stack)

	var je *JsonElement
	var offsetSize int

	for _, rule := range newGrammar {
		lhs := rule.lhs
		expansions := rule.rhs
		for _, expansion := range expansions {
			size := len(expansion)
			if size > stackSize {
				continue
			}
			actual := topNOfStack(stack, size)
			matches := compare(expansion, actual)
			if matches && size > offsetSize {
				je = &JsonElement{
					value:           rule.toJson(stack[len(stack)-size:]...),
					jsonElementType: lhs,
				}
				offsetSize = size
			}
		}
	}

	return je, offsetSize
}

func topNOfStack(stack []*StackElement, count int) []ElementType {
	slice := stack[len(stack)-count:]
	var elements []ElementType

	for _, el := range slice {
		var value = el.value.tokenType
		if el.rule != nil {
			value = el.rule.jsonElementType
		}
		elements = append(elements, value)
	}

	return elements
}

func compare(expansion, actual []ElementType) bool {
	for i := 0; i < len(expansion); i++ {
		if expansion[i] != actual[i] {
			return false
		}
	}
	return true
}

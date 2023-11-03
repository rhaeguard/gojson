package gojson

func action(stack []*stackElement) (*jsonElement, int) {
	stackSize := len(stack)

	var je *jsonElement
	var offsetSize int

	for _, rule := range grammar {
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
				je = &jsonElement{
					value:           rule.toJson(stack[len(stack)-size:]...),
					jsonElementType: lhs,
				}
				offsetSize = size
			}
		}
	}

	return je, offsetSize
}

func topNOfStack(stack []*stackElement, count int) []elementType {
	slice := stack[len(stack)-count:]
	var elements []elementType

	for _, el := range slice {
		var value = el.value.tokenType
		if el.rule != nil {
			value = el.rule.jsonElementType
		}
		elements = append(elements, value)
	}

	return elements
}

func compare(expansion, actual []elementType) bool {
	for i := 0; i < len(expansion); i++ {
		if expansion[i] != actual[i] {
			return false
		}
	}
	return true
}

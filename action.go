package gojson

func action(stack []*stackElement) (*jsonElement, int) {
	stackSize := len(stack)

	var je *jsonElement
	var offset int

	for _, rule := range grammar {
		for _, production := range rule.rhs {
			size := len(production)
			if size > stackSize {
				continue
			}
			actual := topNOfStack(stack, size)
			matches := compare(production, actual)
			if matches && size > offset {
				je = &jsonElement{
					value:           rule.toJson(stack[len(stack)-size:]...),
					jsonElementType: rule.lhs,
				}
				offset = size
			}
		}
	}

	return je, offset
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

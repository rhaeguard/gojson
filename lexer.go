package gojson

import (
	"strings"
)

type token struct {
	value     any
	tokenType elementType
}

var specialSymbols = map[uint8]elementType{
	'{': ltObjectStart,
	'}': ltObjectEnd,
	'[': ltArrayStart,
	']': ltArrayEnd,
	',': ltComma,
	':': ltColon,
	'.': ltFractionSymbol,
}

func isWhitespace(ch uint8) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isDigit(ch uint8) bool {
	return ch >= '0' && ch <= '9'
}

func lex(input string) ([]token, *Error) {
	var tokens []token
	for i := 0; i < len(input); {
		ch := input[i]

		if _, ok := specialSymbols[ch]; ok {
			tokens = append(tokens, token{
				tokenType: specialSymbols[ch],
			})
			i++
		} else if ch == '"' {
			token, offset, err := lexString(input, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += offset
		} else if ch == 't' {
			if "true" == input[i:i+4] {
				tokens = append(tokens, token{
					value:     "true",
					tokenType: ltBoolean,
				})
				i += 4
			} else {
				return nil, newError(i, "unrecognized token")
			}
		} else if ch == 'f' {
			if "false" == input[i:i+5] {
				tokens = append(tokens, token{
					value:     "false",
					tokenType: ltBoolean,
				})
				i += 5
			} else {
				return nil, newError(i, "unrecognized token")
			}
		} else if ch == 'n' {
			if "null" == input[i:i+4] {
				tokens = append(tokens, token{
					tokenType: ltNull,
				})
				i += 4
			} else {
				return nil, newError(i, "unrecognized token")
			}
		} else if isWhitespace(ch) {
			for i < len(input) && isWhitespace(input[i]) {
				i++
			}
		} else if ch == 'e' || ch == 'E' {
			tokens = append(tokens, token{
				tokenType: ltExponent,
			})
			i++
		} else if ch == '+' || ch == '-' {
			tokens = append(tokens, token{
				value:     ch,
				tokenType: ltSign,
			})
			i++
		} else if isDigit(ch) {
			token, offset := lexDigits(input, i)
			tokens = append(tokens, token)
			i += offset
		} else {
			return nil, newError(i, "unrecognized token")
		}
	}
	return tokens, nil
}

func lexDigits(input string, i int) (token, int) {
	var sb strings.Builder
	for i < len(input) && isDigit(input[i]) {
		sb.WriteByte(input[i])
		i++
	}

	return token{
		tokenType: ltDigits,
		value:     sb.String(),
	}, sb.Len()
}

func lexString(input string, i int) (token, int, *Error) {
	i++ // move past the opening quotes
	var sb strings.Builder
	for input[i] != '"' { // TODO: Quote might be escaped
		sb.WriteByte(input[i])
		i++
		if i >= len(input) {
			return token{}, -1, newError(i, "string is not properly closed")
		}
	}

	return token{
			tokenType: ltString,
			value:     sb.String(),
		},
		sb.Len() + 2, // for both the quotes
		nil
}

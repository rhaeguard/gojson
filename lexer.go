package gojson

import (
	"fmt"
	"strings"
)

type SyntaxError struct {
	Pos int
	Msg string
}

func (se *SyntaxError) Error() string {
	if se.Pos == -1 {
		return se.Msg
	}
	return fmt.Sprintf("%s at position %d", se.Msg, se.Pos)
}

func newSyntaxError(Pos int, Msg string) *SyntaxError {
	return &SyntaxError{Pos, Msg}
}

type Token struct {
	value     any
	tokenType ElementType
}

var specialSymbols = map[uint8]ElementType{
	'{': TTObjectStart,
	'}': TTObjectEnd,
	'[': TTArrayStart,
	']': TTArrayEnd,
	',': TTComma,
	':': TTColon,
	'.': TTFractionSymbol,
}

func isWhitespace(ch uint8) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isDigit(ch uint8) bool {
	return ch >= '0' && ch <= '9'
}

func lex(input string) ([]Token, *SyntaxError) {
	var tokens []Token
	for i := 0; i < len(input); {
		ch := input[i]

		switch ch {
		case '{', '}', '[', ']', ',', ':', '.':
			tokens = append(tokens, Token{
				tokenType: specialSymbols[ch],
			})
			i++
		case '"':
			token, offset, err := lexString(input, i)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
			i += offset
		case 't':
			if "true" == input[i:i+4] {
				tokens = append(tokens, Token{
					value:     "true",
					tokenType: TTBoolean,
				})
				i += 4
			} else {
				return nil, newSyntaxError(i, "unrecognized token")
			}
		case 'f':
			if "false" == input[i:i+5] {
				tokens = append(tokens, Token{
					value:     "false",
					tokenType: TTBoolean,
				})
				i += 5
			} else {
				return nil, newSyntaxError(i, "unrecognized token")
			}
		case 'n':
			if "null" == input[i:i+4] {
				tokens = append(tokens, Token{
					tokenType: TTNull,
				})
				i += 4
			} else {
				return nil, newSyntaxError(i, "unrecognized token")
			}
		case ' ', '\t', '\n':
			for i < len(input) && isWhitespace(input[i]) {
				i++
			}
		case 'e', 'E':
			tokens = append(tokens, Token{
				tokenType: TTExponent,
			})
			i++
		case '+', '-':
			tokens = append(tokens, Token{
				value:     ch,
				tokenType: TTSign,
			})
			i++
		default:
			if isDigit(ch) {
				token, offset := lexDigits(input, i)
				tokens = append(tokens, token)
				i += offset
			} else {
				return nil, newSyntaxError(i, "unrecognized token")
			}
		}
	}
	return tokens, nil
}

func lexDigits(input string, i int) (Token, int) {
	var sb strings.Builder
	for i < len(input) && isDigit(input[i]) {
		sb.WriteByte(input[i])
		i++
	}

	return Token{
		tokenType: TTDigits,
		value:     sb.String(),
	}, sb.Len()
}

func lexString(input string, i int) (Token, int, *SyntaxError) {
	i++ // move past the opening quotes
	var sb strings.Builder
	for input[i] != '"' { // TODO: Quote might be escaped
		sb.WriteByte(input[i])
		i++
		if i >= len(input) {
			return Token{}, -1, newSyntaxError(i, "string is not properly closed")
		}
	}

	return Token{
			tokenType: TTString,
			value:     sb.String(),
		},
		sb.Len() + 2, // for both the quotes
		nil
}

package main

import "strings"

type TokenType = string

const (
	TTObjectStart    TokenType = "TT_OBJECT_START"
	TTObjectEnd      TokenType = "TT_OBJECT_END"
	TTArrayStart     TokenType = "TT_ARRAY_START"
	TTArrayEnd       TokenType = "TT_ARRAY_END"
	TTComma          TokenType = "TT_COMMA"
	TTColon          TokenType = "TT_COLON"
	TTFractionSymbol TokenType = "TT_FRACTION_SYMBOL"
	TTBoolean        TokenType = "TT_BOOLEAN"
	TTExponent       TokenType = "TT_EXPONENT"
	TTDigits         TokenType = "TT_DIGITS"
	TTWhitespace     TokenType = "TT_WHITESPACE"
	TTNull           TokenType = "TT_NULL"
	TTSign           TokenType = "TT_SIGN"
	TTString         TokenType = "TT_STRING"
)

type Token struct {
	value     any
	tokenType TokenType
}

var enclosingSymbols = map[uint8]TokenType{
	'{': TTObjectStart,
	'}': TTObjectEnd,
	'[': TTArrayStart,
	']': TTArrayEnd,
}

var specialSymbols = map[uint8]TokenType{
	',': TTComma,
	':': TTColon,
	'.': TTFractionSymbol,
}

func isWhitespace(ch uint8) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func lex(input string) []Token {
	var tokens []Token
	for i := 0; i < len(input); {
		ch := input[i]

		switch ch {
		case '{', '}', '[', ']':
			tokens = append(tokens, Token{
				tokenType: enclosingSymbols[ch],
			})
			i++
		case ',', ':':
			tokens = append(tokens, Token{
				tokenType: specialSymbols[ch],
			})
			i++
		case '"':
			token, offset := lexString(input, i)
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
				panic("invalid syntax")
			}
		case 'f':
			if "false" == input[i:i+5] {
				tokens = append(tokens, Token{
					value:     "false",
					tokenType: TTBoolean,
				})
				i += 5
			} else {
				panic("invalid syntax")
			}
		case 'n':
			if "null" == input[i:i+4] {
				tokens = append(tokens, Token{
					tokenType: TTNull,
				})
				i += 4
			} else {
				panic("invalid syntax")
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
			token, offset := lexDigits(input, i)
			tokens = append(tokens, token)
			i += offset
		}
	}

	return tokens
}

func lexDigits(input string, i int) (Token, int) {
	var sb strings.Builder
	for input[i] >= '0' && input[i] <= '9' {
		sb.WriteByte(input[i])
		i++
	}

	return Token{
		tokenType: TTDigits,
		value:     sb.String(),
	}, sb.Len()
}

func lexString(input string, i int) (Token, int) {
	i++
	var sb strings.Builder
	for input[i] != '"' { // Quote might be escaped
		sb.WriteByte(input[i])
		i++
	}

	return Token{
		tokenType: TTString,
		value:     sb.String(),
	}, sb.Len() + 2 // both quotes
}

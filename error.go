package gojson

import "fmt"

type Error struct {
	Pos int
	Msg string
}

func (se *Error) Error() string {
	if se.Pos == -1 {
		return se.Msg
	}
	return fmt.Sprintf("%s at position %d", se.Msg, se.Pos)
}

func newError(Pos int, Msg string) *Error {
	return &Error{Pos, Msg}
}

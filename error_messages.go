package main

import (
	"fmt"
	"strings"
)

type ErrorType int

const (
	WRONG_ARG_COUNT ErrorType = iota
	WRONG_TYPE
	INT_ERROR
	OVERFLOW_ERROR
)

type ErrorMsg struct {
	Type    ErrorType
	Command string
}

func (e ErrorMsg) Value() Value {
	return Value{
		Type: "error",
		Str:  e.format(),
	}
}

func (e ErrorMsg) format() string {
	switch e.Type {
	case WRONG_ARG_COUNT:
		return fmt.Sprintf(
			"ERR wrong number of arguments for '%s' command",
			strings.ToUpper(e.Command),
		)
	case WRONG_TYPE:
		return "WRONGTYPE Operation against a key holding the wrong kind of value"
	case INT_ERROR:
		return "ERR value is not an integer or out of range"
	case OVERFLOW_ERROR:
		return "ERR increment or decrement would overflow"
	default:
		return "ERR unknown error"
	}
}

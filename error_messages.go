package main

import (
	"fmt"
	"strings"
)

type ErrorType int

const (
	WRONG_ARG_COUNT ErrorType = iota
	WRONG_TYPE
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
	default:
		return "ERR unknown error"
	}
}

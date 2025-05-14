package main

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"COMMAND": func([]Value) Value { return Value{} },
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{t: "string", str: "PONG"}
	}
	return Value{t: "string", str: args[0].bulk}
}

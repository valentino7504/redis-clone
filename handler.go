package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

var (
	dict      = map[string]string{}
	setsMutex = sync.RWMutex{}
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{t: "string", str: "PONG"}
	}
	return Value{t: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) < 2 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	value := args[1].bulk
	setsMutex.Lock()
	dict[key] = value
	setsMutex.Unlock()
	return Value{t: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) < 2 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'GET' command"}
	}
	key := args[0].bulk
	setsMutex.Lock()
	value, ok := dict[key]
	if !ok {
		return Value{t: "null"}
	}
	return Value{t: "bulk", bulk: value}
}

package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"HSET":    hset,
	"GET":     get,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var (
	setDict    = map[string]string{}
	setsMutex  = sync.RWMutex{}
	hsetDict   = map[string]map[string]string{}
	hsetsMutex = sync.RWMutex{}
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{t: "string", str: "PONG"}
	}
	return Value{t: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].bulk
	value := args[1].bulk
	setsMutex.Lock()
	setDict[key] = value
	setsMutex.Unlock()
	return Value{t: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'GET' command"}
	}
	key := args[0].bulk
	setsMutex.RLock()
	value, ok := setDict[key]
	if !ok {
		return Value{t: "null"}
	}
	setsMutex.RUnlock()
	return Value{t: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'HSET' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk
	hsetsMutex.Lock()
	if _, ok := hsetDict[hash]; !ok {
		hsetDict[hash] = map[string]string{}
	}
	hsetDict[hash][key] = value
	hsetsMutex.Unlock()
	return Value{t: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'HGET' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	hsetsMutex.RLock()
	val, ok := hsetDict[hash][key]
	if !ok {
		return Value{t: "null"}
	}
	hsetsMutex.RUnlock()
	return Value{t: "bulk", bulk: val}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{t: "error", str: "ERR wrong number of arguments for 'HGET' command"}
	}
	hash := args[0].bulk
	hsetsMutex.RLock()
	if _, ok := hsetDict[hash]; !ok {
		return Value{t: "array", array: []Value{}}
	}
	values := make([]Value, 0)
	for k, v := range hsetDict[hash] {
		values = append(values, Value{t: "bulk", bulk: k}, Value{t: "bulk", bulk: v})
	}
	hsetsMutex.RUnlock()
	return Value{t: "array", array: values}
}

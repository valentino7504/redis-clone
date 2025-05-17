package main

import "sync"

type RedisType int

const (
	STRINGTYPE RedisType = iota
	HASHTYPE
	// refactored like this to allow for newer types as I go
)

type RedisValue struct {
	Type      RedisType
	StringVal string
	HashVal   map[string]string
}

var Handlers = map[string]func([]Value) Value{
	// generic handlers
	"PING":     ping,
	"DEL":      del,
	"TYPE":     typeCheck,
	"EXISTS":   exists,
	"FLUSHALL": flushall,
	// string handlers
	"SET": set,
	"GET": get,
	// hash handlers
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"HDEL":    hdel,
}

var (
	store      = map[string]RedisValue{}
	storeMutex = sync.RWMutex{}
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Type: "string", Str: "PONG"}
	}
	return Value{Type: "string", Str: args[0].Bulk}
}

func exists(args []Value) Value {
	count := 0
	storeMutex.RLock()
	for _, key := range args {
		if _, ok := store[key.Bulk]; ok {
			count++
		}
	}
	storeMutex.RUnlock()
	return Value{Type: "integer", Int: count}
}

func typeCheck(args []Value) Value {
	if len(args) != 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "TYPE"}.Value()
	}
	key := args[0]
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	val, ok := store[key.Bulk]
	if !ok {
		return Value{Type: "string", Str: "none"}
	}
	switch val.Type {
	case STRINGTYPE:
		return Value{Type: "string", Str: "string"}
	case HASHTYPE:
		return Value{Type: "string", Str: "hash"}
	default:
		return Value{Type: "string", Str: "none"}
	}
}

func del(args []Value) Value {
	if len(args) < 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "DEL"}.Value()
	}
	storeMutex.Lock()
	count := 0
	for _, key := range args {
		if _, ok := store[key.Bulk]; ok {
			count++
			delete(store, key.Bulk)
		}
	}
	storeMutex.Unlock()
	return Value{Type: "integer", Int: count}
}

func flushall(args []Value) Value {
	if len(args) != 0 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "DEL"}.Value()
	}
	storeMutex.Lock()
	store = make(map[string]RedisValue)
	storeMutex.Unlock()
	return Value{Type: "string", Str: "OK"}
}

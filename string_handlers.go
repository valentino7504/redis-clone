package main

func set(args []Value) Value {
	if len(args) != 2 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "SET"}.Value()
	}
	key := args[0].Bulk
	value := args[1].Bulk
	storeMutex.Lock()
	store[key] = RedisValue{Type: STRINGTYPE, StringVal: value}
	storeMutex.Unlock()
	return Value{Type: "string", Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "GET"}.Value()
	}
	key := args[0].Bulk
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	value, ok := store[key]
	if !ok {
		return Value{Type: "null"}
	}
	if value.Type != STRINGTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	}
	return Value{Type: "bulk", Bulk: value.StringVal}
}

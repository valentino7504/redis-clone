package main

func hset(args []Value) Value {
	if len(args) != 3 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "HSET"}.Value()
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk
	storeMutex.Lock()
	if value, ok := store[hash]; !ok || value.Type != HASHTYPE {
		store[hash] = RedisValue{Type: HASHTYPE, HashVal: make(map[string]string)}
	}
	store[hash].HashVal[key] = value
	storeMutex.Unlock()
	return Value{Type: "string", Str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "HGET"}.Value()
	}
	hash := args[0].Bulk
	key := args[1].Bulk
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	storedVal, ok := store[hash]
	if !ok {
		return Value{Type: "null"}
	}
	if storedVal.Type != HASHTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	}
	val, ok := storedVal.HashVal[key]
	if !ok {
		return Value{Type: "null"}
	}
	return Value{Type: "bulk", Bulk: val}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "HGETALL"}.Value()
	}
	hash := args[0].Bulk
	storeMutex.RLock()
	defer storeMutex.RUnlock()
	if _, ok := store[hash]; !ok {
		return Value{Type: "array", Array: []Value{}}
	}
	if store[hash].Type != HASHTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	}
	values := make([]Value, 0)
	for k, v := range store[hash].HashVal {
		values = append(values, Value{Type: "bulk", Bulk: k}, Value{Type: "bulk", Bulk: v})
	}
	return Value{Type: "array", Array: values}
}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "HDEL"}.Value()
	}
	hash := args[0].Bulk
	storeMutex.Lock()
	defer storeMutex.Unlock()
	storedVal, ok := store[hash]
	var count int64
	if !ok {
		return Value{Type: "integer", Int: count}
	}
	if storedVal.Type != HASHTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	}
	dict := storedVal.HashVal
	for _, key := range args[1:] {
		if _, ok := dict[key.Bulk]; ok {
			count++
			delete(store[hash].HashVal, key.Bulk)
		}
	}
	return Value{Type: "integer", Int: count}
}

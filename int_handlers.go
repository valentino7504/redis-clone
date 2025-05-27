package main

import (
	"math"
	"strconv"
)

// Increments an integer by 1. Checks if the value is a string type and is an
// integer that is parseable, and if yes, then adds 1 to it as long as it is within
// the overflow
func incr(args []Value) Value {
	if len(args) != 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "INCR"}.Value()
	}
	key := args[0].Bulk
	storeMutex.Lock()
	defer storeMutex.Unlock()
	str, ok := store[key]
	if ok && str.Type != STRINGTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	} else if !ok {
		store[key] = RedisValue{Type: STRINGTYPE, StringVal: "0"}
		str = store[key]
	}
	i, err := strconv.ParseInt(str.StringVal, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	if i == math.MaxInt64 {
		return ErrorMsg{Type: OVERFLOW_ERROR}.Value()
	}
	i++
	newVal := strconv.FormatInt(i, 10)
	store[key] = RedisValue{Type: STRINGTYPE, StringVal: newVal}
	return Value{Type: "integer", Int: i}
}

func decr(args []Value) Value {
	if len(args) != 1 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "DECR"}.Value()
	}
	key := args[0].Bulk
	storeMutex.Lock()
	defer storeMutex.Unlock()
	str, ok := store[key]
	if ok && str.Type != STRINGTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	} else if !ok {
		store[key] = RedisValue{Type: STRINGTYPE, StringVal: "0"}
		str = store[key]
	}
	i, err := strconv.ParseInt(str.StringVal, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	if i == math.MinInt64 {
		return ErrorMsg{Type: OVERFLOW_ERROR}.Value()
	}
	i--
	newVal := strconv.FormatInt(i, 10)
	store[key] = RedisValue{Type: STRINGTYPE, StringVal: newVal}
	return Value{Type: "integer", Int: i}
}

func incrby(args []Value) Value {
	if len(args) != 2 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "INCRBY"}.Value()
	}
	key, increment := args[0].Bulk, args[1].Bulk
	incrementInt, err := strconv.ParseInt(increment, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	storeMutex.Lock()
	defer storeMutex.Unlock()
	str, ok := store[key]
	if ok && str.Type != STRINGTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	} else if !ok {
		store[key] = RedisValue{Type: STRINGTYPE, StringVal: "0"}
		str = store[key]
	}
	i, err := strconv.ParseInt(str.StringVal, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	if math.MaxInt64-incrementInt <= i {
		return ErrorMsg{Type: OVERFLOW_ERROR}.Value()
	}
	i += incrementInt
	newVal := strconv.FormatInt(i, 10)
	store[key] = RedisValue{Type: STRINGTYPE, StringVal: newVal}
	return Value{Type: "integer", Int: i}
}

func decrby(args []Value) Value {
	if len(args) != 2 {
		return ErrorMsg{Type: WRONG_ARG_COUNT, Command: "DECRBY"}.Value()
	}
	key, decrement := args[0].Bulk, args[1].Bulk
	decrementInt, err := strconv.ParseInt(decrement, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	storeMutex.Lock()
	defer storeMutex.Unlock()
	str, ok := store[key]
	if ok && str.Type != STRINGTYPE {
		return ErrorMsg{Type: WRONG_TYPE}.Value()
	} else if !ok {
		store[key] = RedisValue{Type: STRINGTYPE, StringVal: "0"}
		str = store[key]
	}
	i, err := strconv.ParseInt(str.StringVal, 10, 64)
	if err != nil {
		return ErrorMsg{Type: INT_ERROR}.Value()
	}
	if math.MinInt64+decrementInt >= i {
		return ErrorMsg{Type: OVERFLOW_ERROR}.Value()
	}
	i -= decrementInt
	newVal := strconv.FormatInt(i, 10)
	store[key] = RedisValue{Type: STRINGTYPE, StringVal: newVal}
	return Value{Type: "integer", Int: i}
}

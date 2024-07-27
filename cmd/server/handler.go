package server

import "sync"

var Handlers = map[string]func(args []Value) Value{
	"PING": ping,
	"INFO": info,
	"ECHO": echo,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func info(args []Value) Value {
	return Value{typ: "string", str: "Welcome to Redis!"}
}

func echo(args []Value) Value {
	resultStr := ""

	for i, arg := range args {
		if i > 0 {
			resultStr += " "
		}

		resultStr += arg.bulk
	}

	return Value{typ: "string", str: resultStr}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "SET command must have exactly 2 arguments"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "GET command must have exactly 1 argument"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	value, ok := SETs[key]
	SETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "HSET command must have exactly 3 arguments"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "HGET command must have exactly 2 arguments"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.Lock()
	value, ok := HSETs[hash][key]
	HSETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

package server

var Handlers = map[string]func(args []Value) Value{
	"PING": ping,
	"INFO": info,
	"ECHO": echo,
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

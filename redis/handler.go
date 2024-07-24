package redis

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
	return Value{typ: "string", str: args[0].bulk}
}

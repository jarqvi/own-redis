package redis

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var (
	listen = flag.String("listen", ":6379", "address to listen on")
)

func Server() (err error) {
	flag.Parse()

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %v", *listen, err)
	}

	defer close(l, &err, "failed to close listener")

	log.Printf("server listening on %s", *listen)

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}

		log.Printf("accepted connection from %s", c.RemoteAddr())

		go func(c net.Conn) {
			defer close(c, &err, "failed to close connection")

			err := cmd(c)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}(c)
	}
}

func cmd(c net.Conn) error {
	for {
		resp := NewResp(c)
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				log.Printf("client disconnected: %s", c.RemoteAddr())
				return nil
			}

			return fmt.Errorf("failed to read data from connection: %v", err)
		}

		writer := NewWriter(c)

		if value.typ != "array" {
			err = writer.Write(Value{typ: "error", str: "expected array"})
			if err != nil {
				return fmt.Errorf("failed to write data to connection: %v", err)
			}

			continue
		}

		if len(value.array) == 0 {
			err = writer.Write(Value{typ: "error", str: "expected array with at least one element"})
			if err != nil {
				return fmt.Errorf("failed to write data to connection: %v", err)
			}

			continue
		}

		log.Printf("received: %v", value)

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[cmd]
		if !ok {
			err = writer.Write(Value{typ: "error", str: "unknown command"})
			if err != nil {
				return fmt.Errorf("failed to write data to connection: %v", err)
			}

			continue
		}

		err = writer.Write(handler(args))
		if err != nil {
			return fmt.Errorf("failed to write data to connection: %v", err)
		}
	}
}

func close(c io.Closer, errP *error, msg string) {
	err := c.Close()
	if err != nil {
		*errP = fmt.Errorf("%s: %v", msg, err)
	}
}

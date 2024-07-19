package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var (
	listen = flag.String("listen", ":6379", "address to listen on")
)

func main() {
	flag.Parse()

	err := server()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func server() (err error) {
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

		log.Printf("received: %v", value)

		writer := NewWriter(c)
		err = writer.Write(Value{typ: "string", str: "OK"})
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

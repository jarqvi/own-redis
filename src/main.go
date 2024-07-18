package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	listen = flag.String("listen", ":6379", "address to listen on")
)

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() (err error) {
	l, err := net.Listen("tcp", *listen)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %v", *listen, err)
	}

	defer closeResource(l, &err, "failed to close listener")

	log.Printf("server listening on %s", *listen)

	for {
		c, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %v", err)
		}

		log.Printf("accepted connection from %s", c.RemoteAddr())

		go func(c net.Conn) {
			defer closeResource(c, &err, "failed to close connection")

			err := cmd(c)
			if err != nil {
				log.Printf("error: %v", err)
			}
		}(c)
	}
}

func cmd(c net.Conn) error {
	buf := make([]byte, 128)

	for {
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("client disconnected: %s", c.RemoteAddr())
				return nil
			}
			return fmt.Errorf("failed to read data from connection: %v", err)
		}

		cmd := strings.TrimSpace(string(buf[:n]))

		log.Printf("Read command: %s", cmd)

		if strings.ToLower(cmd) == "ping" {
			_, err = c.Write([]byte("+PONG\r\n"))
			if err != nil {
				return fmt.Errorf("failed to write data to connection: %v", err)
			}
		} else {
			_, err = c.Write([]byte("(error) unknown command '" + cmd + "'\r\n"))
			if err != nil {
				return fmt.Errorf("failed to write data to connection: %v", err)
			}
		}
	}
}

func closeResource(c io.Closer, errP *error, msg string) {
	err := c.Close()
	if err != nil {
		*errP = fmt.Errorf("%s: %v", msg, err)
	}
}

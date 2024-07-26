package main

import (
	"fmt"
	"os"

	"github.com/jarqvi/own-redis/cmd/server"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

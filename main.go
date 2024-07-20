package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jarqvi/own-redis/redis"
)

func main() {
	flag.Parse()

	err := redis.Server()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

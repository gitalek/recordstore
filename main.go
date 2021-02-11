package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("error while connecting: %#v\n", err)
	}
	defer conn.Close()

	_, err = conn.Do(
		"HMSET", "album:2",
		"title", "Electric Ladyland",
		"artist", "Jimi Hendrix",
		"price", 4.95,
		"likes", 8,
	)
	if err != nil {
		log.Fatalf("error while doing command: %#v\n", err)
	}
	fmt.Printf("Electric LadyLand added!\n")
}

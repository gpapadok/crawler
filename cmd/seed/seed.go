package main

import (
	"fmt"
	"crawler/broker"

	"github.com/joho/godotenv"
)

// TODO: Maybe move to DB
var seedURL string = "http://example.com"

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Create broker connection
	conn, err := broker.Connect()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to broker.")

	// Create broker channel
	ch, err := broker.CreateChannel(conn)
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	fmt.Println("Successfully created broker channel.")

	// Create broker Queue
	queueName := "urls" // TODO: Move to const
	q, err := broker.CreateQueue(ch, queueName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created broker queue.")

	err = broker.Publish(ch, q, seedURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("Published body", seedURL)
}

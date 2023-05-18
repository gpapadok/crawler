package main

import (
	"fmt"
	"web-crawler/broker"
	"web-crawler/database"
	"web-crawler/scraper"
)

func main() {
	// Create database connection
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Successfully connected to database.")

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
	queueName := "urls"
	q, err := broker.CreateQueue(ch, queueName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created broker queue.")

	for {
		msgs, err := broker.Consume(ch, q)
		if err != nil {
			panic(err)
		}

		for msg := range msgs {
			fmt.Println("Consumed: ", string(msg.Body))
			err = scraper.Scrape(string(msg.Body), db, ch, q)
			if err != nil {
				panic(err)
			}
		}
	}
}

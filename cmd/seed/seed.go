package main

import (
	"crawler/broker"
	"crawler/database"
	"log"

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
	log.Println("Successfully connected to broker.")

	// Create broker channel
	ch, err := broker.CreateChannel(conn)
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	log.Println("Successfully created broker channel.")

	// Create broker Queue
	queueName := "urls" // TODO: Move to const
	q, err := broker.CreateQueue(ch, queueName)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully created broker queue.")
	}

	db, err := database.Connect(database.DBInfo)
	if err != nil {
		log.Println(err)

		err = broker.Publish(ch, q, seedURL)
		if err != nil {
			panic(err)
		}
		log.Println("Published body", seedURL)
		return
	}

	urls, err := database.GetUnscraped(db)
	if err != nil {
		log.Println(err)
	}

	if err != nil || len(urls) == 0 {
		err = broker.Publish(ch, q, seedURL)
		if err != nil {
			panic(err)
		}
		log.Println("Published body", seedURL)
		return
	}

	for _, url := range urls {
		err = broker.Publish(ch, q, url)
		if err != nil {
			panic(err)
		}
	}
	log.Println("Published unscraped urls from database.")
}

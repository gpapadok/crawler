package main

import (
	"crawler/broker"
	"crawler/database"
	"crawler/scraper"
	"log"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const N_WORKERS = 4

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	crawl()
}

func crawl() {
	// Create database connection
	db, err := database.Connect(database.DBInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Println("Successfully connected to database.")

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
	queueName := "urls"
	q, err := broker.CreateQueue(ch, queueName)
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created broker queue.")

	links := make(chan scraper.Link)
	scrape := func(url string) error {
		return scraper.Scrape(url, links)
	}

	go scraper.StoreAndPublish(links, scraper.Connections{db, ch, q})

	deliveries, err := broker.Consume(ch, q)
	if err != nil {
		panic(err)
	}

	sema := make(chan struct{}, N_WORKERS)

	for msg := range deliveries {
		sema <- struct{}{}
		go handleMsg(scrape, msg, sema)
	}
}

func handleMsg(scrape func(string) error, msg amqp.Delivery, sema <-chan struct{}) {
	defer func() { <-sema }()
	if err := scrape(string(msg.Body)); err != nil {
		log.Println(err)
		return
	}

	if err := msg.Ack(false); err != nil {
		log.Println(err)
	}
}

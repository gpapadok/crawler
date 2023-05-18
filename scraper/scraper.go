package scraper

import (
	"database/sql"
	"web-crawler/broker"
	"web-crawler/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MOCK PARSER
func parse(url string) ([]string, error) {
	return []string{"url1", "url2", "url3"}, nil
}

func Scrape(url string, db *sql.DB, channel *amqp.Channel, queue *amqp.Queue) error {
	links, err := parse(url)
	if err != nil {
		return err
	}

	for _, link := range links {
		err = broker.Publish(channel, queue, link)
		if err != nil {
			return err
		}

		err = database.InsertURL(db, link, url)
		if err != nil {
			return err
		}
	}
	return nil
}

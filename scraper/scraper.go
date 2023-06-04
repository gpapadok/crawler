package scraper

import (
	"database/sql"
	"fmt"
	"net/http"
	"web-crawler/broker"
	"web-crawler/database"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
)

func Scrape(url string, db *sql.DB, channel *amqp.Channel, queue *amqp.Queue) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil // TODO: Custom error
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Scraping links: ", url)
	links := visit([]string{}, doc)

	for _, link := range links {
		link = buildLink(url, link)

		isVisited, err := database.IsVisited(db, link)
		if err != nil {
			return err
		}
		if isVisited {
			continue
		}

		if err = broker.Publish(channel, queue, link); err != nil {
			return err
		}

		if err = database.InsertURL(db, link, url); err != nil {
			return err
		}
	}
	return nil
}

func visit(links []string, node *html.Node) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			// Filter out links with protocols other than http and #fragment links
			if a.Key == "href" && validLink(a.Val) {
				links = append(links, a.Val)
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

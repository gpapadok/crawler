package scraper

import (
	"crawler/broker"
	"crawler/database"
	"database/sql"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/html"
)

type Connections struct {
	Db      *sql.DB
	Channel *amqp.Channel
	Queue   *amqp.Queue
}

type Link struct {
	url    string
	parent string
}

func Scrape(url string, links chan Link) error {
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

	log.Println("Scraping links: ", url)
	traverseForLinks(links, doc, url)

	return nil
}

func StoreAndPublish(links chan Link, conns Connections) {
	for link := range links {
		url := buildLink(link.parent, link.url)

		isVisited, err := database.IsVisited(conns.Db, url)
		if err != nil {
			log.Println(err)
			continue
		}
		if isVisited {
			continue
		}

		if err = database.InsertURL(conns.Db, url, link.parent); err != nil {
			log.Println(err)
			continue
		}
		if err = broker.Publish(conns.Channel, conns.Queue, url); err != nil {
			log.Println(err)
			continue
		}
	}
}

func traverseForLinks(links chan Link, node *html.Node, parent string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			// Filter out links with protocols other than http and #fragment links
			if a.Key == "href" && validLink(a.Val) {
				links <- Link{a.Val, parent}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseForLinks(links, c, parent)
	}
}

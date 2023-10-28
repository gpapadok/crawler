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

func Scrape(url string, links chan<- Link) error {
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

func linksBuilder(links <-chan Link, builtLinks chan<- Link) {
	for link := range links {
		builtLinks <- Link{buildLink(link.parent, link.url), link.parent}
	}
}

func StoreAndPublish(links <-chan Link, conns Connections) {
	builtLinks := make(chan Link)
	go linksBuilder(links, builtLinks)

	for link := range builtLinks {
		isVisited, err := database.IsVisited(conns.Db, link.url)
		if err != nil {
			log.Println(err)
			continue
		}
		if isVisited {
			continue
		}

		if err = database.InsertURL(conns.Db, link.url, link.parent); err != nil {
			log.Println(err)
			continue
		}
		if err = broker.Publish(conns.Channel, conns.Queue, link.url); err != nil {
			log.Println(err)
		}
	}
}

func traverseForLinks(links chan<- Link, node *html.Node, parent string) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			// Filter out non http and #fragment links
			if a.Key == "href" && validLink(a.Val) {
				links <- Link{a.Val, parent}
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		traverseForLinks(links, c, parent)
	}
}

# Crawler

A distributed web crawler written in Go.
RabbitMQ is used as the message broker and supervisor and PostgreSQL for persistent storage.

Each node consumes a url link from the broker, fetches the HTML and parses it for links.
It then publishes each unvisited link in the broker and stores it in the database.

### Usage
Credentials for the broker and database should be stored in a `.env` file. See below for a sample configuration.

To run locally you need docker installed.

Build seed and node container images with `make build`

Then initialze PostgreSQL and Rabbitmq containers and seed them by running `make init`

Finally to start a crawler node `make run`

#### dotenv

```
DB_HOST=crawler-db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=crawler

TEST_DB_HOST=crawler-db
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=crawler_test

BROKER_HOST=crawler-mq
BROKER_PORT=5672
BROKER_USER=guest
BROKER_PASSWORD=guest
```

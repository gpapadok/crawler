# Crawler

A distributed web crawler written in Go.
RabbitMQ is used as the message broker and supervisor and PostgreSQL for persistent storage.

Each node consumes a url link from the broker, fetches the HTML and parses it for links.
It then publishes each unvisited link in the broker and stores it in the database.

### Usage
Credentials for the broker and database should be stored in a `.env` file. See below for a sample configuration to run locally.

Initialze PostgreSQL and Rabbitmq containers by running `start_containers.sh`

Build commands:
```
go build cmd/seed/seed.go
go build cmd/node/node.go
```

Seed the broker with a URL `./seed`

Then start running nodes with `./node`

#### dotenv

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=crawler

TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=crawler_test

BROKER_HOST=localhost
BROKER_PORT=5672
BROKER_USER=guest
BROKER_PASSWORD=guest
```

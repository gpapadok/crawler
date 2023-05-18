package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "crawler"
)

func psqlInfo() (psqlInfo string) {
	psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return
}

func Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertURL(db *sql.DB, url string, parent string) error {
	statement := `insert into links (parent, url)
	    values ($1, $2) on conflict do nothing`

	_, err := db.Exec(statement, parent, url)

	return err
}

func GetVisited(db *sql.DB) ([]string, error) {
	query := "select distinct parent from links"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var link string
	var links []string

	for rows.Next() {
		rows.Scan(&link)
		links = append(links, link)
	}

	return links, nil
}

func IsVisited(db *sql.DB, link string) (bool, error) {
	query := "select $1 in (select url from links)"

	row := db.QueryRow(query, link)

	var isVisited bool

	err := row.Scan(&isVisited)
	if err != nil {
		return false, err
	}

	return isVisited, nil
}

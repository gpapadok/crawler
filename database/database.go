package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func DBInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
}

func Connect(buildDbInfo func() string) (*sql.DB, error) {
	db, err := sql.Open("postgres", buildDbInfo())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
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

func IsVisited(db *sql.DB, link string) (bool, error) {
	query := "select $1 in (select parent from links)"

	var isVisited bool
	err := db.QueryRow(query, link).Scan(&isVisited)

	return isVisited, err
}

func GetUnscraped(db *sql.DB) ([]string, error) {
	query := "select url from links where url not in (select parent from links)"

	var urls []string
	var url string
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&url)
		if err != nil {
			panic(err)
		}
		urls = append(urls, url)
	}

	return urls, err
}

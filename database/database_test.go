package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func testDBInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))
}

func connectTest() (*sql.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", testDBInfo())

	return db, err
}

func clearTestDB(db *sql.DB) error {
	cmd := "truncate links"
	_, err := db.Exec(cmd)
	return err
}

func countURL(db *sql.DB) (int, error) {
	query := "select count(*) from links"

	row := db.QueryRow(query)

	var count int

	err := row.Scan(&count)

	return count, err
}

func TestInsertURL(t *testing.T) {
	db, err := connectTest()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	countBefore, err := countURL(db)
	if err != nil {
		panic(err)
	}

	err = InsertURL(db, "http://example.com/route", "http://example.com")
	if err != nil {
		panic(err)
	}

	countAfter, err := countURL(db)
	if err != nil {
		panic(err)
	}

	if countAfter-countBefore != 1 {
		t.Error("Insert failed", countBefore, countAfter, countAfter-countBefore)
	}

	clearTestDB(db)
}

func TestIsVisited(t *testing.T) {
	db, err := connectTest()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	parent := "http://example.com/route"
	link := "http://example.com"

	if err = InsertURL(db, link, parent); err != nil {
		panic(err)
	}

	parentVisited, err := IsVisited(db, parent)
	if err != nil {
		panic(err)
	}

	if !parentVisited {
		t.Error("Parent not visited")
	}

	linkVisited, err := IsVisited(db, link)
	if err != nil {
		panic(err)
	}

	if linkVisited {
		t.Error("Link visited")
	}

	clearTestDB(db)
}

func TestGetUnscraped(t *testing.T) {
	db, err := connectTest()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	links := []struct {
		parent string
		url    string
	}{
		{"http://example.com", "http://example.com/route"},
		{"http://example.com", "http://example.com/hello"},
		{"http://example.com", "http://example.com/alekos"},
		{"http://example.com/alekos", "http://example.com/alekos/index.html"},
	}
	nChildren := 3

	for _, link := range links {
		if err = InsertURL(db, link.url, link.parent); err != nil {
			panic(err)
		}
	}

	urls, err := GetUnscraped(db)
	if err != nil {
		panic(err)
	}

	if len(urls) != nChildren {
		t.Errorf(`# unscraped = %d, want %d`, len(urls), nChildren)
	}
	clearTestDB(db)
}

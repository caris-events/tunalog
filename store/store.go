package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Instance *Store

type Store struct {
	db *sql.DB
}

func init() {
	db, err := sql.Open("sqlite3", "file:db.sqlite?cache=shared")
	if err != nil {
		log.Fatalln("open: %w", err)
	}

	Instance = &Store{db}
	if err := createUserTable(); err != nil {
		log.Fatalln("exec: %w", err)
	}
	if err := createPostTable(); err != nil {
		log.Fatalln("exec: %w", err)
	}
	if err := createNavigationTable(); err != nil {
		log.Fatalln("exec: %w", err)
	}
	if err := createTagTable(); err != nil {
		log.Fatalln("exec: %w", err)
	}
	if err := createPostTagTable(); err != nil {
		log.Fatalln("exec: %w", err)
	}
	//defer db.Close()
}

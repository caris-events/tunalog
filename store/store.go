package store

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite", "file:db.sqlite?cache=shared")
	if err != nil {
		log.Fatalln(err)
	}
	if err := createUserTable(); err != nil {
		log.Fatalln(err)
	}
	if err := createPostTable(); err != nil {
		log.Fatalln(err)
	}
	if err := createNavigationTable(); err != nil {
		log.Fatalln(err)
	}
	if err := createTagTable(); err != nil {
		log.Fatalln(err)
	}
	if err := createPostTagTable(); err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			if err := ClearExpiredTrashPosts(); err != nil {
				log.Println(err)
			}
			<-time.After(24 * time.Hour)
		}
	}()
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

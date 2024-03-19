package store

import "github.com/caris-events/tunalog/entity"

type Post struct {
	ID          string
	Title       string
	Slug        string
	Excerpt     string
	CoverPath   string
	AuthorID    string
	Password    string
	Visibility  entity.Visibility
	Content     string
	PinnedAt    int64
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
}

func createPostTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
		id             TEXT NOT NULL PRIMARY KEY,
		title          TEXT NOT NULL,
		slug           TEXT NOT NULL,
		excerpt        TEXT NOT NULL,
		author_id      TEXT NOT NULL,
		cover_path     TEXT NOT NULL,
		password       TEXT NOT NULL,
		visibility     TEXT NOT NULL,
		content        TEXT NOT NULL,
		pinned_at      INTEGER NOT NULL,
		published_at   INTEGER NOT NULL,
		created_at     INTEGER NOT NULL,
		updated_at     INTEGER NOT NULL
	)
`)
	return err
}

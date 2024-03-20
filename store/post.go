package store

import (
	"context"

	"github.com/caris-events/tunalog/entity"
)

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

func (s *Store) CountPostsByUser(c context.Context, uid string) (int, error) {
	row := s.db.QueryRow("SELECT COUNT(*) FROM posts WHERE author_id = ?", uid)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// TransferPosts
func (s *Store) TransferPosts(c context.Context, fromUserID, toUserID string) error {
	_, err := s.db.Exec("UPDATE posts SET author_id = ? WHERE author_id = ?", toUserID, fromUserID)
	if err != nil {
		return err
	}
	return nil
}

// DeletePostsByUser
func (s *Store) DeletePostsByUser(c context.Context, id string) error {
	_, err := s.db.Exec("DELETE FROM post_tags WHERE post_id IN (SELECT id FROM posts WHERE author_id = ?)", id) // implicit linking table
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM posts WHERE author_id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

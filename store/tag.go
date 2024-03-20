package store

import (
	"context"
	"database/sql"

	"github.com/caris-events/tunalog/entity"
)

type Tag struct {
	ID          string
	Slug        string
	Name        string
	Description string
	CreatedAt   int64
}

func createTagTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS tags (
		id          TEXT NOT NULL PRIMARY KEY,
		slug        TEXT NOT NULL,
		name        TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at  INTEGER NOT NULL
	)
`)
	return err
}

func createPostTagTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS post_tags (
		tag_id  TEXT NOT NULL,
		post_id TEXT NOT NULL,
		PRIMARY KEY (tag_id, post_id)
	)
`)
	return err
}

// ListTags
func (s *Store) ListTags(c context.Context, offset, limit int, keyword string) ([]*entity.TagR, error) {
	var rows *sql.Rows
	var err error

	if keyword != "" {
		rows, err = s.db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id WHERE name LIKE ? GROUP BY t.id, t.slug, t.name, t.description ORDER BY t.created_at DESC LIMIT ?, ?`, "%"+keyword+"%", offset, limit)
	} else {
		rows, err = s.db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id GROUP BY t.id, t.slug, t.name, t.description ORDER BY t.created_at DESC LIMIT ?, ?`, offset, limit)

	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.TagR
	for rows.Next() {
		var t entity.TagR
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt, &t.PostCount); err != nil {
			return nil, err
		}

		tags = append(tags, &t)
	}
	return tags, nil
}

// CountTags
func (s *Store) CountTags(c context.Context, keyword string) (int, error) {
	var count int
	if keyword != "" {
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM tags WHERE name LIKE ?`, "%"+keyword+"%").Scan(&count); err != nil {
			return 0, err
		}
	} else {
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

// CreateTag
func (s *Store) CreateTag(c context.Context, t *entity.TagW) error {
	_, err := s.db.Exec(`INSERT INTO tags (id, slug, name, description, created_at) VALUES (?, ?, ?, ?, ?)`, t.ID, t.Slug, t.Name, t.Description, t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

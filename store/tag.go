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

// UpdateTag
func (s *Store) UpdateTag(c context.Context, t *entity.TagW) error {
	_, err := s.db.Exec(`UPDATE tags SET slug = ?, name = ?, description = ?, created_at = ? WHERE id = ?`, t.Slug, t.Name, t.Description, t.CreatedAt, t.ID)
	if err != nil {
		return err
	}
	return nil
}

// ListTagsByPost
func (s *Store) ListTagsByPost(c context.Context, id string) ([]*entity.TagR, error) {
	rows, err := s.db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.TagR
	for rows.Next() {
		var t entity.TagR
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

// GetTag
func (s *Store) GetTag(c context.Context, id string) (*entity.TagR, error) {
	var t entity.TagR
	if err := s.db.QueryRow(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id WHERE t.id = ? GROUP BY t.id, t.slug, t.name, t.description`, id).Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt, &t.PostCount); err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTag
func (s *Store) DeleteTag(id string) error {
	_, err := s.db.Exec(`DELETE FROM post_tags WHERE tag_id = ?`, id) // implicit linking table
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM tags WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return nil

}

// GetTagsByName
func (s *Store) GetTagsByName(names []string) ([]*entity.TagR, error) {
	var args []interface{}

	q := "SELECT id, slug, name, description, created_at FROM tags WHERE name IN ("
	for i := 0; i < len(names); i++ {
		q += "?"
		args = append(args, names[i])
		if i < len(names)-1 {
			q += ","
		}
	}
	q += ")"

	rows, err := s.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*entity.TagR
	for rows.Next() {
		var t entity.TagR
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

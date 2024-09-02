package store

import (
	"database/sql"

	"github.com/caris-events/tunalog/entity"
)

func createTagTable() error {
	_, err := db.Exec(`
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
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS post_tags (
		tag_id  TEXT NOT NULL,
		post_id TEXT NOT NULL,
		PRIMARY KEY (tag_id, post_id)
	)
`)
	return err
}

func ListTags(offset, limit int, keyword string) ([]*entity.TagR, error) {
	var rows *sql.Rows
	var err error

	if keyword != "" {
		rows, err = db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id WHERE name LIKE ? GROUP BY t.id, t.slug, t.name, t.description ORDER BY t.created_at DESC LIMIT ?, ?`, "%"+keyword+"%", offset, limit)
	} else {
		rows, err = db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id GROUP BY t.id, t.slug, t.name, t.description ORDER BY t.created_at DESC LIMIT ?, ?`, offset, limit)

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

func CountTags(keyword string) (int, error) {
	var count int
	if keyword != "" {
		if err := db.QueryRow(`SELECT COUNT(*) FROM tags WHERE name LIKE ?`, "%"+keyword+"%").Scan(&count); err != nil {
			return 0, err
		}
	} else {
		if err := db.QueryRow(`SELECT COUNT(*) FROM tags`).Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}

func CreateTag(t *entity.TagW) error {
	if _, err := db.Exec(`INSERT INTO tags (id, slug, name, description, created_at) VALUES (?, ?, ?, ?, ?)`, t.ID, t.Slug, t.Name, t.Description, t.CreatedAt); err != nil {
		return err
	}
	return nil
}

func GetTag(id string) (*entity.TagR, error) {
	var t entity.TagR
	if err := db.QueryRow(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t LEFT JOIN post_tags pt ON t.id = pt.tag_id WHERE t.id = ? GROUP BY t.id, t.slug, t.name, t.description`, id).Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt, &t.PostCount); err != nil {
		return nil, err
	}
	return &t, nil
}

func ListMostUsedTags() ([]*entity.TagR, error) {
	rows, err := db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at, COUNT(pt.post_id) AS post_count FROM tags t JOIN post_tags pt ON t.id = pt.tag_id GROUP BY t.id, t.slug, t.name, t.description ORDER BY post_count DESC LIMIT 10`)
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

func UpdateTag(t *entity.TagW) error {
	if _, err := db.Exec(`UPDATE tags SET slug = ?, name = ?, description = ?, created_at = ? WHERE id = ?`, t.Slug, t.Name, t.Description, t.CreatedAt, t.ID); err != nil {
		return err
	}
	return nil
}

func DeleteTag(id string) error {
	if _, err := db.Exec(`DELETE FROM post_tags WHERE tag_id = ?`, id); err != nil {
		return err
	}
	if _, err := db.Exec(`DELETE FROM tags WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func GetTagsByName(names []string) ([]*entity.TagR, error) {
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

	rows, err := db.Query(q, args...)
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

func GetTagBySlug(slug string) (*entity.TagR, error) {
	var t entity.TagR
	if err := db.QueryRow(`SELECT id, slug, name, description, created_at FROM tags WHERE slug = ?`, slug).Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func ListTagsByPost(id string) ([]*entity.TagR, error) {
	rows, err := db.Query(`SELECT t.id, t.slug, t.name, t.description, t.created_at FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = ?`, id)
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

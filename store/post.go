package store

import (
	"time"

	"github.com/caris-events/tunalog/entity"
	"github.com/caris-events/tunalog/system"
)

func createPostTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
		id             TEXT NOT NULL PRIMARY KEY,
		title          TEXT NOT NULL,
		slug           TEXT NOT NULL,
		excerpt        TEXT NOT NULL,
		author_id      TEXT NOT NULL,
		password       TEXT NOT NULL,
		visibility     TEXT NOT NULL,
		content        TEXT NOT NULL,
		pinned_at      INTEGER NOT NULL,
		published_at   INTEGER NOT NULL,
		created_at     INTEGER NOT NULL,
		updated_at     INTEGER NOT NULL,
		trashed_at     INTEGER NOT NULL
	)
`)
	return err
}

func CountPostsByUser(uid string) (int, error) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM posts WHERE author_id = ?", uid).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func TransferPosts(fromUID, toUID string) error {
	if _, err := db.Exec("UPDATE posts SET author_id = ? WHERE author_id = ?", toUID, fromUID); err != nil {
		return err
	}
	return nil
}

func DeletePostsByUser(uid string) error {
	if _, err := db.Exec("DELETE FROM post_tags WHERE post_id IN (SELECT id FROM posts WHERE author_id = ?)", uid); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM posts WHERE author_id = ?", uid); err != nil {
		return err
	}
	return nil
}

func CreatePost(p *entity.PostW) error {
	if _, err := db.Exec("INSERT INTO posts (id, title, slug, excerpt, author_id, password, visibility, content, published_at, created_at, updated_at, pinned_at, trashed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", p.ID, p.Title, p.Slug, p.Excerpt, p.AuthorID, p.Password, p.Visibility, p.Content, p.PublishedAt, p.CreatedAt, p.UpdatedAt, p.PinnedAt, p.TrashedAt); err != nil {
		return err
	}
	for _, tagID := range p.TagIDs {
		if _, err := db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", p.ID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func TrashPost(id string) error {
	if _, err := db.Exec("UPDATE posts SET trashed_at = ? WHERE id = ?", time.Now().Unix(), id); err != nil {
		return err
	}
	return nil
}

func UntrashPost(id string) error {
	if _, err := db.Exec("UPDATE posts SET trashed_at = ? WHERE id = ?", 0, id); err != nil {
		return err
	}
	return nil
}

func DeletePost(id string) error {
	if _, err := db.Exec("DELETE FROM post_tags WHERE post_id = ?", id); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM posts WHERE id = ? AND trashed_at != ?", id, 0); err != nil {
		return err
	}
	return nil
}

func ClearTrashPosts() error {
	if _, err := db.Exec("DELETE FROM post_tags WHERE post_id IN (SELECT id FROM posts WHERE trashed_at != ?)", 0); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM posts WHERE trashed_at != ?", 0); err != nil {
		return err
	}
	return nil
}

func ClearExpiredTrashPosts() error {
	if _, err := db.Exec("DELETE FROM post_tags WHERE post_id IN (SELECT id FROM posts WHERE trashed_at != ? AND trashed_at < ?)", 0, time.Now().Unix()-86400*30); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM posts WHERE trashed_at != ? AND trashed_at < ?", 0, time.Now().Unix()-86400*30); err != nil {
		return err
	}
	return nil
}

func GetPreviousPost(id string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := db.QueryRow(`
        SELECT p.id, p.title, p.slug, p.excerpt, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, u.nickname, u.email, u.bio, u.created_at
        FROM posts p
        JOIN users u ON p.author_id = u.id
        WHERE p.published_at < ?
        AND (p.published_at < (SELECT published_at FROM posts WHERE id = ?) OR (p.published_at = (SELECT published_at FROM posts WHERE id = ?) AND p.created_at < (SELECT created_at FROM posts WHERE id = ?)))
        AND (p.visibility = ? OR p.visibility = ?)
        AND p.trashed_at = 0
        ORDER BY p.published_at DESC, p.created_at DESC
        LIMIT 1`,
		time.Now().Unix(), id, id, id, "public", "password").Scan(
		&p.ID, &p.Title, &p.Slug, &p.OriginalExcerpt, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.Author.Nickname, &p.Author.Email, &p.Author.Bio, &p.Author.CreatedAt); err != nil {
		return nil, err
	}
	tags, err := ListTagsByPost(p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

func GetNextPost(id string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := db.QueryRow(`
        SELECT p.id, p.title, p.slug, p.excerpt, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, u.nickname, u.email, u.bio, u.created_at
        FROM posts p
        JOIN users u ON p.author_id = u.id
        WHERE p.published_at < ?
        AND (p.published_at > (SELECT published_at FROM posts WHERE id = ?) OR (p.published_at = (SELECT published_at FROM posts WHERE id = ?) AND p.created_at > (SELECT created_at FROM posts WHERE id = ?)))
        AND (p.visibility = ? OR p.visibility = ?)
        AND p.trashed_at = 0
        ORDER BY p.published_at ASC, p.created_at ASC
        LIMIT 1`,
		time.Now().Unix(), id, id, id, "public", "password").Scan(
		&p.ID, &p.Title, &p.Slug, &p.OriginalExcerpt, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.Author.Nickname, &p.Author.Email, &p.Author.Bio, &p.Author.CreatedAt); err != nil {
		return nil, err
	}
	tags, err := ListTagsByPost(p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

func PtrBool(b bool) *bool {
	return &b
}

type ListPostsQuery struct {
	AuthorID       string
	TagID          string
	Title          string
	Query          string
	Visibilities   []entity.Visibility
	IsPublished    *bool
	IsTrashed      *bool
	PublishedYear  string
	PublishedMonth string
	PublishedDay   string
	PublishedDate  string
	Offset         int
	Limit          int
}

func (q *ListPostsQuery) Build() (query string, args []any) {
	if q.TagID != "" {
		query += " JOIN post_tags pt ON p.id = pt.post_id"
	}
	query += " WHERE 1 = 1"

	if q.TagID != "" {
		query += " AND pt.tag_id = ?"
		args = append(args, q.TagID)
	}
	if q.AuthorID != "" {
		query += " AND p.author_id = ?"
		args = append(args, q.AuthorID)
	}
	if q.Title != "" {
		query += " AND p.title LIKE ?"
		args = append(args, "%"+q.Title+"%")
	}
	if q.Query != "" {
		query += " AND (p.title LIKE ? OR p.content LIKE ?)"
		args = append(args, "%"+q.Query+"%", "%"+q.Query+"%")
	}
	if q.IsPublished != nil {
		if *q.IsPublished {
			query += " AND p.published_at <= ?"
			args = append(args, time.Now().Unix())
		} else {
			query += " AND p.published_at > ?"
			args = append(args, time.Now().Unix())
		}
	}
	if q.IsTrashed != nil {
		if *q.IsTrashed {
			query += " AND p.trashed_at != 0"
		} else {
			query += " AND p.trashed_at = 0"
		}
	}
	if len(q.Visibilities) > 0 {
		query += " AND p.visibility IN ("
		for i, v := range q.Visibilities {
			if i > 0 {
				query += ", "
			}
			query += "?"
			args = append(args, v)
		}
		query += ")"
	}
	if q.PublishedYear != "" {
		query += " AND strftime('%Y', datetime(published_at + ?, 'unixepoch')) = ?"
		args = append(args, system.Config.Timezone, q.PublishedYear)
	}
	if q.PublishedMonth != "" {
		query += " AND strftime('%m', datetime(published_at + ?, 'unixepoch')) = ?"
		args = append(args, system.Config.Timezone, q.PublishedMonth)
	}
	if q.PublishedDay != "" {
		query += " AND strftime('%d', datetime(published_at + ?, 'unixepoch')) = ?"
		args = append(args, system.Config.Timezone, q.PublishedDay)
	}
	if q.PublishedDate != "" {
		query += " AND strftime('%Y-%m', datetime(published_at + ?, 'unixepoch')) = ?"
		args = append(args, system.Config.Timezone, q.PublishedDate)
	}
	return
}

func CountPosts(q *ListPostsQuery) (int, error) {
	var args []any
	query := "SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id"

	qQuery, qArgs := q.Build()
	query += qQuery
	args = append(args, qArgs...)

	var count int
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func CountPostsByType() (*entity.PostCount, error) {
	rows, err := db.Query("SELECT visibility, COUNT(*) FROM posts WHERE trashed_at = 0 GROUP BY visibility")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts entity.PostCount
	for rows.Next() {
		var v string
		var c int
		if err := rows.Scan(&v, &c); err != nil {
			return nil, err
		}
		switch v {
		case "public":
			counts.Public = c
		case "private":
			counts.Private = c
		case "password":
			counts.Password = c
		case "draft":
			counts.Draft = c
		}
	}

	row := db.QueryRow("SELECT COUNT(*) FROM posts WHERE trashed_at != 0")
	if err := row.Scan(&counts.Trash); err != nil {
		return nil, err
	}

	counts.All = counts.Public + counts.Private + counts.Password + counts.Draft + counts.Trash
	counts.NonTrash = counts.All - counts.Trash
	return &counts, nil
}

func ListPostDates() (data []string, err error) {
	rows, err := db.Query("SELECT strftime('%Y-%m', datetime(published_at, 'unixepoch')) FROM posts GROUP BY strftime('%Y-%m', datetime(published_at, 'unixepoch'))")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return data, nil
}

func ListPosts(q *ListPostsQuery) ([]*entity.PostR, error) {
	var args []any
	query := "SELECT p.id, p.title, p.slug, p.excerpt, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.trashed_at, u.id, u.nickname, u.email, u.bio, u.created_at FROM posts p JOIN users u ON p.author_id = u.id"

	qQuery, qArgs := q.Build()
	query += qQuery
	args = append(args, qArgs...)
	query += " ORDER BY p.pinned_at DESC, p.published_at DESC"

	if q.Limit > 0 && q.Offset >= 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, q.Limit, q.Offset)
	} else if q.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, q.Limit)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.PostR
	for rows.Next() {
		var p entity.PostR
		p.Author = entity.UserR{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.OriginalExcerpt, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.TrashedAt, &p.Author.ID, &p.Author.Nickname, &p.Author.Email, &p.Author.Bio, &p.Author.CreatedAt); err != nil {
			return nil, err
		}
		tags, err := ListTagsByPost(p.ID)
		if err != nil {
			return nil, err
		}
		p.Tags = tags
		posts = append(posts, &p)
	}
	return posts, nil
}

func GetPost(id string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := db.QueryRow("SELECT p.id, p.title, p.slug, p.excerpt, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.trashed_at, u.id, u.nickname, u.email, u.bio, u.created_at FROM posts p JOIN users u ON p.author_id = u.id WHERE p.id = ? AND p.trashed_at = 0", id).Scan(&p.ID, &p.Title, &p.Slug, &p.OriginalExcerpt, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.TrashedAt, &p.Author.ID, &p.Author.Nickname, &p.Author.Email, &p.Author.Bio, &p.Author.CreatedAt); err != nil {
		return nil, err
	}
	tags, err := ListTagsByPost(p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

func UpdatePost(p *entity.PostW) error {
	if _, err := db.Exec("UPDATE posts SET title = ?, slug = ?, excerpt = ?, author_id = ?, password = ?, visibility = ?, content = ?, published_at = ?, created_at = ?, updated_at = ?, pinned_at = ? WHERE id = ?", p.Title, p.Slug, p.Excerpt, p.AuthorID, p.Password, p.Visibility, p.Content, p.PublishedAt, p.CreatedAt, p.UpdatedAt, p.PinnedAt, p.ID); err != nil {
		return err
	}
	if _, err := db.Exec("DELETE FROM post_tags WHERE post_id = ?", p.ID); err != nil {
		return err
	}
	for _, tagID := range p.TagIDs {
		if _, err := db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", p.ID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func GetPostBySlug(slug string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := db.QueryRow("SELECT p.id, p.title, p.slug, p.excerpt, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.trashed_at, u.id, u.nickname, u.email, u.bio, u.created_at FROM posts p JOIN users u ON p.author_id = u.id WHERE p.slug = ? AND p.trashed_at = 0", slug).Scan(&p.ID, &p.Title, &p.Slug, &p.OriginalExcerpt, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.TrashedAt, &p.Author.ID, &p.Author.Nickname, &p.Author.Email, &p.Author.Bio, &p.Author.CreatedAt); err != nil {
		return nil, err
	}
	tags, err := ListTagsByPost(p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

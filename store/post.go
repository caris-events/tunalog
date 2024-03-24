package store

import (
	"context"
	"time"

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
	IsDraft     bool
	IsHidden    bool
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
		updated_at     INTEGER NOT NULL,
		is_draft       INTEGER NOT NULL,
		is_hidden      INTEGER NOT NULL
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

type ListPostsOrderBy string

const (
	ListPostsOrderByCreatedAtDesc ListPostsOrderBy = ""
	ListPostsOrderByCreatedAtAsc  ListPostsOrderBy = "created_at_asc"
)

type ListPostsQuery struct {
	AuthorID       string
	TagID          string
	Title          string
	Query          string
	Visibility     entity.Visibility
	HasPassword    *bool
	HasPublished   *bool
	PublishedYear  string
	PublishedMonth string
	PublishedDay   string
	PublishedDate  string
	OrderBy        ListPostsOrderBy
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
	if q.HasPassword != nil {
		if *q.HasPassword {
			query += " AND p.password != ''"
		} else {
			query += " AND p.password = ''"
		}
	}
	if q.HasPublished != nil {
		if *q.HasPublished {
			query += " AND p.published_at <= ?"
			args = append(args, time.Now().Unix())
		} else {
			query += " AND p.published_at > ?"
			args = append(args, time.Now().Unix())
		}
	}
	if q.Visibility != entity.VisibilityUnknown {
		query += " AND p.visibility = ?"
		args = append(args, q.Visibility)
	}
	if q.PublishedYear != "" {
		query += " AND strftime('%Y', datetime(published_at, 'unixepoch')) = ?"
		args = append(args, q.PublishedYear)
	}
	if q.PublishedMonth != "" {
		query += " AND strftime('%m', datetime(published_at, 'unixepoch')) = ?"
		args = append(args, q.PublishedMonth)
	}
	if q.PublishedDay != "" {
		query += " AND strftime('%d', datetime(published_at, 'unixepoch')) = ?"
		args = append(args, q.PublishedDay)
	}
	if q.PublishedDate != "" {
		query += " AND strftime('%Y-%m', datetime(published_at, 'unixepoch')) = ?"
		args = append(args, q.PublishedDate)
	}
	return
}

func (s *Store) ListPosts(c context.Context, q *ListPostsQuery) ([]*entity.PostR, error) {
	var args []any
	query := "SELECT p.id, p.title, p.slug, p.excerpt, p.cover_path, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.is_hidden, p.is_draft, u.username, u.nickname, u.bio, u.created_at, u.suspended_at FROM posts p JOIN users u ON p.author_id = u.id"

	qQuery, qArgs := q.Build()
	query += qQuery
	args = append(args, qArgs...)

	switch q.OrderBy {
	case ListPostsOrderByCreatedAtDesc:
		query += " ORDER BY p.created_at DESC"
	case ListPostsOrderByCreatedAtAsc:
		query += " ORDER BY p.created_at ASC"
	}

	if q.Limit > 0 && q.Offset >= 0 {
		query += " LIMIT ? OFFSET ?"
		args = append(args, q.Limit, q.Offset)
	} else if q.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, q.Limit)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entity.PostR
	for rows.Next() {
		var p entity.PostR
		p.Author = entity.UserR{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.CoverPath, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.IsHidden, &p.IsDraft, &p.Author.Username, &p.Author.Nickname, &p.Author.Bio, &p.Author.CreatedAt, &p.Author.SuspendedAt); err != nil {
			return nil, err
		}
		tags, err := s.ListTagsByPost(c, p.ID)
		if err != nil {
			return nil, err
		}
		p.Tags = tags
		posts = append(posts, &p)
	}
	return posts, nil
}

func (s *Store) CountPosts(c context.Context, q *ListPostsQuery) (int, error) {
	var args []any
	query := "SELECT COUNT(*) FROM posts p JOIN users u ON p.author_id = u.id"

	qQuery, qArgs := q.Build()
	query += qQuery
	args = append(args, qArgs...)

	var count int

	if err := s.db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Store) ListPostDates(c context.Context) (data []string, err error) {
	rows, err := s.db.Query("SELECT strftime('%Y-%m', datetime(published_at, 'unixepoch')) FROM posts GROUP BY strftime('%Y-%m', datetime(published_at, 'unixepoch'))")
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

// CreatePost
func (s *Store) CreatePost(c context.Context, p *entity.PostW) error {
	_, err := s.db.Exec("INSERT INTO posts (id, title, slug, excerpt, cover_path, author_id, password, visibility, content, published_at, created_at, updated_at, pinned_at, is_draft, is_hidden) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", p.ID, p.Title, p.Slug, p.Excerpt, p.CoverPath, p.AuthorID, p.Password, p.Visibility, p.Content, p.PublishedAt, p.CreatedAt, p.UpdatedAt, p.PinnedAt, p.IsDraft, p.IsHidden)
	if err != nil {
		return err
	}
	for _, tagID := range p.TagIDs {
		_, err := s.db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", p.ID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetPostBySlug
func (s *Store) GetPostBySlug(c context.Context, slug string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := s.db.QueryRow("SELECT p.id, p.title, p.slug, p.excerpt, p.cover_path, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.is_hidden, p.is_draft, u.username, u.nickname, u.bio, u.created_at, u.suspended_at FROM posts p JOIN users u ON p.author_id = u.id WHERE p.slug = ?", slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.CoverPath, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.IsHidden, &p.IsDraft, &p.Author.Username, &p.Author.Nickname, &p.Author.Bio, &p.Author.CreatedAt, &p.Author.SuspendedAt); err != nil {
		return nil, err
	}
	tags, err := s.ListTagsByPost(c, p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

// GetPost
func (s *Store) GetPost(c context.Context, id string) (*entity.PostR, error) {
	var p entity.PostR
	p.Author = entity.UserR{}
	if err := s.db.QueryRow("SELECT p.id, p.title, p.slug, p.excerpt, p.cover_path, p.author_id, p.password, p.visibility, p.content, p.published_at, p.created_at, p.updated_at, p.pinned_at, p.is_hidden, p.is_draft, u.username, u.nickname, u.bio, u.created_at, u.suspended_at FROM posts p JOIN users u ON p.author_id = u.id WHERE p.id = ?", id).Scan(&p.ID, &p.Title, &p.Slug, &p.Excerpt, &p.CoverPath, &p.AuthorID, &p.Password, &p.Visibility, &p.Content, &p.PublishedAt, &p.CreatedAt, &p.UpdatedAt, &p.PinnedAt, &p.IsHidden, &p.IsDraft, &p.Author.Username, &p.Author.Nickname, &p.Author.Bio, &p.Author.CreatedAt, &p.Author.SuspendedAt); err != nil {
		return nil, err
	}
	tags, err := s.ListTagsByPost(c, p.ID)
	if err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

// UpdatePost
func (s *Store) UpdatePost(c context.Context, p *entity.PostW) error {
	_, err := s.db.Exec("UPDATE posts SET title = ?, slug = ?, excerpt = ?, cover_path = ?, author_id = ?, password = ?, visibility = ?, content = ?, published_at = ?, created_at = ?, updated_at = ?, pinned_at = ?, is_hidden = ?, is_draft = ? WHERE id = ?", p.Title, p.Slug, p.Excerpt, p.CoverPath, p.AuthorID, p.Password, p.Visibility, p.Content, p.PublishedAt, p.CreatedAt, p.UpdatedAt, p.PinnedAt, p.IsHidden, p.IsDraft, p.ID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM post_tags WHERE post_id = ?", p.ID)
	if err != nil {
		return err
	}
	for _, tagID := range p.TagIDs {
		_, err := s.db.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", p.ID, tagID)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeletePost
func (s *Store) DeletePost(c context.Context, id string) error {
	_, err := s.db.Exec("DELETE FROM post_tags WHERE post_id = ?", id) // implicit linking table
	if err != nil {
		return err
	}
	_, err = s.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

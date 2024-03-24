package entity

import (
	"strings"
	"text/template"
	"time"
)

type Visibility string

const (
	VisibilityUnknown  Visibility = ""
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityPassword Visibility = "password"
)

type PostW struct {
	ID          string
	Title       string
	Slug        string
	Excerpt     string
	CoverPath   string
	AuthorID    string
	Password    string
	Visibility  Visibility
	Content     string
	PinnedAt    int64
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
	IsDraft     bool
	IsHidden    bool

	TagIDs []string
}

type PostR struct {
	ID          string
	Title       string
	Slug        string
	Excerpt     string
	CoverPath   string
	AuthorID    string
	Password    string
	Visibility  Visibility
	Content     string
	PinnedAt    int64
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
	IsDraft     bool
	IsHidden    bool

	Author UserR
	Tags   []*TagR
}

func (p *PostR) TagsStr() string {
	var tags []string
	for _, t := range p.Tags {
		tags = append(tags, t.Name)
	}
	return template.HTMLEscapeString(strings.Join(tags, ","))
}

func (p *PostR) PublishedAtDatetime() string {
	return time.Unix(p.PublishedAt, 0).Format("2006-01-02 15:04:05")
}

func (p *PostR) IsUnpublished() bool {
	return time.Now().Unix() < p.PublishedAt
}

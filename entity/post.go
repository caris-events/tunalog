package entity

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	stripmd "github.com/writeas/go-strip-markdown"
)

type Visibility string

const (
	VisibilityUnknown  Visibility = ""
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityPassword Visibility = "password"
	VisibilityDraft    Visibility = "draft"
	VisibilityTrash    Visibility = "trash" // trash doesn't really exists, it's an abstraction
)

type PostCount struct {
	All      int
	NonTrash int
	Public   int
	Private  int
	Password int
	Draft    int
	Trash    int
}

type PostW struct {
	ID          string
	Title       string
	Slug        string
	Excerpt     string
	AuthorID    string
	Password    string
	Visibility  Visibility
	Content     string
	PinnedAt    int64
	PublishedAt int64
	CreatedAt   int64
	UpdatedAt   int64
	TrashedAt   int64

	TagIDs []string
}

type PostR struct {
	ID              string
	Title           string
	Slug            string
	OriginalExcerpt string
	AuthorID        string
	Password        string
	Visibility      Visibility
	Content         string
	PinnedAt        int64
	PublishedAt     int64
	CreatedAt       int64
	UpdatedAt       int64
	TrashedAt       int64

	Author UserR
	Tags   []*TagR
}

func (p *PostR) TagsStr() string {
	return strings.Join(p.TagNames(), ",")
}

func (p *PostR) TagNames() []string {
	var tags []string
	for _, tag := range p.Tags {
		tags = append(tags, tag.Name)
	}
	if len(tags) == 0 {
		return make([]string, 0)
	}
	return tags
}

func (p *PostR) PublishedDate() string {
	return time.Unix(p.PublishedAt, 0).Format("2006-01-02")
}

func (p *PostR) PublishedAtDatetime() string {
	return time.Unix(p.PublishedAt, 0).Format("2006-01-02 03:04 PM")
}

func (p *PostR) PublishedAtISO() string {
	return time.Unix(p.PublishedAt, 0).Format("2006-01-02T15:04")
}

func (p *PostR) PublishedYear() string {
	return time.Unix(p.PublishedAt, 0).Format("2006")
}

func (p *PostR) PublishedMonth() string {
	return time.Unix(p.PublishedAt, 0).Format("01")
}

func (p *PostR) PublishedDay() string {
	return time.Unix(p.PublishedAt, 0).Format("02")
}

func (p *PostR) IsPublished() bool {
	return time.Now().Unix() >= p.PublishedAt
}

func (p *PostR) Cover() string {
	if _, err := os.Stat(fmt.Sprintf("data/uploads/covers/%s.jpg", p.ID)); os.IsNotExist(err) {
		return ""
	}
	return fmt.Sprintf("uploads/covers/%s.jpg", p.ID)
}

func (p *PostR) Excerpt() string {
	if p.OriginalExcerpt != "" {
		return p.OriginalExcerpt
	}
	content := stripmd.Strip(p.Content)
	if utf8.RuneCountInString(content) > 200 {
		return string([]rune(content)[:200]) + "..."
	}
	return content
}

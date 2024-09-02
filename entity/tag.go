package entity

type TagW struct {
	ID          string
	Slug        string
	Name        string
	Description string
	CreatedAt   int64
}

type TagR struct {
	ID          string
	Slug        string
	Name        string
	Description string
	CreatedAt   int64
	PostCount   int
}

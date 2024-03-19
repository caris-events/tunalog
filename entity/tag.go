package entity

type TagW struct {
	ID          string
	Slug        string
	Name        string
	Description string
}

type TagR struct {
	ID          string
	Slug        string
	Name        string
	Description string
	PostCount   int
}

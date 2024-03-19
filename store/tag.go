package store

type Tag struct {
	ID          string
	Slug        string
	Name        string
	Description string
}

func createTagTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS tags (
		id          TEXT NOT NULL PRIMARY KEY,
		slug        TEXT NOT NULL,
		name        TEXT NOT NULL,
		description TEXT NOT NULL
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

package store

type Navigation struct {
	ID       string
	URL      string
	Name     string
	Sequence int
}

func createNavigationTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS navigations (
		id         TEXT NOT NULL PRIMARY KEY,
		url        TEXT NOT NULL,
		name       TEXT NOT NULL,
		sequence   INTEGER NOT NULL
	)
`)
	return err
}

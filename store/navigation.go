package store

import (
	"github.com/caris-events/tunalog/entity"
)

func createNavigationTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS navigations (
		id         TEXT NOT NULL PRIMARY KEY,
		url        TEXT NOT NULL,
		name       TEXT NOT NULL,
		sequence   INTEGER NOT NULL
	)
`)
	return err
}

func ListNavigations() ([]*entity.NavigationR, error) {
	rows, err := db.Query(`SELECT id, url, name, sequence FROM navigations ORDER BY sequence ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var navigations []*entity.NavigationR
	for rows.Next() {
		var n entity.NavigationR
		if err := rows.Scan(&n.ID, &n.URL, &n.Name, &n.Sequence); err != nil {
			return nil, err
		}
		navigations = append(navigations, &n)
	}
	return navigations, nil
}

func ClearNavigations() error {
	if _, err := db.Exec(`DELETE FROM navigations`); err != nil {
		return err
	}
	return nil
}

func CreateNavigation(n *entity.NavigationW) error {
	if _, err := db.Exec(`INSERT INTO navigations (id, url, name, sequence) VALUES (?, ?, ?, ?)`, n.ID, n.URL, n.Name, n.Sequence); err != nil {
		return err
	}
	return nil
}

package store

import (
	"context"

	"github.com/caris-events/tunalog/entity"
)

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

// ListNavigations
func (s *Store) ListNavigations(c context.Context) ([]*entity.NavigationR, error) {
	rows, err := s.db.Query(`SELECT id, url, name, sequence FROM navigations ORDER BY sequence ASC`)
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

// ClearNavigations
func (s *Store) ClearNavigations(c context.Context) error {
	_, err := s.db.Exec(`DELETE FROM navigations`)
	if err != nil {
		return err
	}
	return nil
}

// CreateNavigation
func (s *Store) CreateNavigation(c context.Context, n *entity.NavigationW) error {
	_, err := s.db.Exec(`INSERT INTO navigations (id, url, name, sequence) VALUES (?, ?, ?, ?)`, n.ID, n.URL, n.Name, n.Sequence)
	if err != nil {
		return err
	}
	return nil
}

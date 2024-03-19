package store

import (
	"context"

	"github.com/caris-events/tunalog/entity"
)

type User struct {
	ID          string
	Username    string
	Password    string
	Nickname    string
	Bio         string
	AvatarPath  string
	CreatedAt   int64
	SuspendedAt int64
}

func createUserTable() error {
	_, err := Instance.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id             TEXT NOT NULL PRIMARY KEY,
		nickname       TEXT NOT NULL UNIQUE,
		username       TEXT NOT NULL UNIQUE,
		password       TEXT NOT NULL,
		bio            TEXT NOT NULL,
		avatar_path    TEXT NOT NULL,
		suspended_at   INTEGER NOT NULL,
		created_at     INTEGER NOT NULL
	)
`)
	return err
}

// CreateUser
func (s *Store) CreateUser(c context.Context, u *entity.UserW) error {
	_, err := s.db.Exec(`INSERT INTO users (id, username, password, nickname, bio, avatar_path, created_at, suspended_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, u.ID, u.Username, u.Password, u.Nickname, u.Bio, u.AvatarPath, u.CreatedAt, u.SuspendedAt)
	if err != nil {
		return err
	}
	return nil
}

// ListUsers
func (s *Store) ListUsers(c context.Context) ([]*entity.UserR, error) {
	rows, err := s.db.Query(`SELECT id, username, nickname, bio, avatar_path, created_at, suspended_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserR
	for rows.Next() {
		var u entity.UserR
		if err := rows.Scan(&u.ID, &u.Username, &u.Nickname, &u.Bio, &u.AvatarPath, &u.CreatedAt, &u.SuspendedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

// GetUser
func (s *Store) GetUser(c context.Context, id string) (*entity.UserR, error) {
	row := s.db.QueryRow(`SELECT id, username, password, nickname, bio, avatar_path, created_at, suspended_at FROM users WHERE id = ?`, id)
	var u entity.UserR
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Nickname, &u.Bio, &u.AvatarPath, &u.CreatedAt, &u.SuspendedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

package store

import (
	"github.com/caris-events/tunalog/entity"
)

func createUserTable() error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id             TEXT NOT NULL PRIMARY KEY,
		email 		   TEXT NOT NULL UNIQUE,
		nickname       TEXT NOT NULL UNIQUE,
		password       TEXT NOT NULL,
		bio            TEXT NOT NULL,
		created_at     INTEGER NOT NULL
	)
`)
	return err
}

func CreateUser(u *entity.UserW) error {
	if _, err := db.Exec(`INSERT INTO users (id, email, nickname, password, bio, created_at) VALUES (?, ?, ?, ?, ?, ?)`, u.ID, u.Email, u.Nickname, u.Password, u.Bio, u.CreatedAt); err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*entity.UserR, error) {
	var u entity.UserR
	if err := db.QueryRow(`SELECT id, email, nickname, password, bio, created_at FROM users WHERE email = ?`, email).Scan(&u.ID, &u.Email, &u.Nickname, &u.Password, &u.Bio, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func ListUsers() ([]*entity.UserR, error) {
	rows, err := db.Query(`SELECT id, email, nickname, password, bio, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserR
	for rows.Next() {
		var u entity.UserR
		if err := rows.Scan(&u.ID, &u.Email, &u.Nickname, &u.Password, &u.Bio, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func GetUser(id string) (*entity.UserR, error) {
	var u entity.UserR
	if err := db.QueryRow(`SELECT id, email, nickname, password, bio, created_at FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Email, &u.Nickname, &u.Password, &u.Bio, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUser(id, nickname, bio, email string) error {
	if _, err := db.Exec(`UPDATE users SET nickname = ?, bio = ?, email = ? WHERE id = ?`, nickname, bio, email, id); err != nil {
		return err
	}
	return nil
}

func UpdateUserPassword(id, password string) error {
	if _, err := db.Exec(`UPDATE users SET password = ? WHERE id = ?`, password, id); err != nil {
		return err
	}
	return nil
}

func DeleteUser(id string) error {
	if _, err := db.Exec(`DELETE FROM users WHERE id = ?`, id); err != nil {
		return err
	}
	return nil
}

func UserExists(id string) (bool, error) {
	var exists bool
	if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

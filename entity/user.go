package entity

import (
	"crypto/md5"
	"fmt"
)

type UserW struct {
	ID        string
	Email     string
	Password  string
	Nickname  string
	Bio       string
	CreatedAt int64
}

type UserR struct {
	ID        string
	Email     string
	Password  string
	Nickname  string
	Bio       string
	CreatedAt int64
}

func (u *UserR) Gravatar() string {
	data := []byte(u.Email)
	return fmt.Sprintf("http://www.gravatar.com/avatar/%x", md5.Sum(data))
}

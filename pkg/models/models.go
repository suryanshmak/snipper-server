package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrDuplicateEmail     = errors.New("models: email already exists")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

type User struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

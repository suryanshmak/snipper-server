package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"suryanshmak.net/snippetBox/pkg/models"
)

type UserModel struct {
	DB *pgx.Conn
}

func (u *UserModel) Insert(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	exec := `INSERT INTO users (name, email, password, created) 
			 VALUES ($1, $2, $3, NOW())`

	_, err = u.DB.Exec(context.Background(), exec, name, email, hash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (u *UserModel) Authenticate(email, pasword string) error {
	stmt := `SELECT password FROM users WHERE email = $1`

	var hashedPassword []byte
	err := u.DB.QueryRow(context.Background(), stmt, email).Scan(&hashedPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.NoDataFound {
				return models.ErrInvalidCredentials
			}
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(pasword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return models.ErrInvalidCredentials
	} else if err != nil {
		return err
	}

	return nil
}

func (u *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

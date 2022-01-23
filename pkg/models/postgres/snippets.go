package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"
	"suryanshmak.net/snippetBox/pkg/models"
)

type SnippetModel struct {
	DB *pgx.Conn
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires) 
	VALUES ($1, $2, NOW(), now() + $3)`
	_, err := m.DB.Exec(context.Background(), stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	var id int
	err = m.DB.QueryRow(context.Background(), "SELECT id FROM snippets ORDER BY created DESC LIMIT 1").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > NOW() AND id = $1`
	s := &models.Snippet{}
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > now() ORDER BY created DESC LIMIT 10`
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

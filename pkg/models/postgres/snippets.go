package postgres

import (
	"database/sql"
	"errors"

	"chrisshyi.net/mini_url/pkg/models"
)

// MiniURLModel is used for database operations
type MiniURLModel struct {
	DB *sql.DB
}

// Insert : insert a snippet
func (m *MiniURLModel) Insert(URL string) (int, error) {
	// TODO: Accommodate durations other than days?
	stmt := `INSERT INTO mini_urls (url, visits)
    VALUES($1, 1) RETURNING id`

	var newMiniURLID int
	err := m.DB.QueryRow(stmt, URL).Scan(&newMiniURLID)
	if err != nil {
		return 0, err
	}

	return newMiniURLID, nil
}

// Get : retrieve a snippet
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
			 WHERE expires > NOW() AND id = $1`
	s := &models.Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

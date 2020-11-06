package postgres

import (
	"database/sql"
	"fmt"

	"chrisshyi.net/mini_url/pkg/models"
)

// MiniURLModel is used for database operations
type MiniURLModel struct {
	DB *sql.DB
}

// GetByID retrieves a MiniURL by its ID field
func (m *MiniURLModel) GetByID(ID int) (*models.MiniURL, error) {
	fmt.Printf("Getting miniURL by ID %d\n", ID)
	stmt := `SELECT id, url, visits FROM mini_urls WHERE id = $1`
	miniURL := &models.MiniURL{}
	err := m.DB.QueryRow(stmt, ID).Scan(&miniURL.ID, &miniURL.URL, &miniURL.Visits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return miniURL, nil
}

// GetByURL retrieves a MiniURL by its url field
func (m *MiniURLModel) GetByURL(URL string) (*models.MiniURL, error) {
	stmt := `SELECT id, url, visits FROM mini_urls WHERE url = $1`
	miniURL := &models.MiniURL{}
	err := m.DB.QueryRow(stmt, URL).Scan(&miniURL.ID, &miniURL.URL, &miniURL.Visits)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return miniURL, nil
}

// Insert : insert a new URL
func (m *MiniURLModel) Insert(URL string) (int, error) {
	stmt := `INSERT INTO mini_urls (url, visits)
    VALUES($1, 1) RETURNING id`

	var newMiniURLID int
	err := m.DB.QueryRow(stmt, URL).Scan(&newMiniURLID)
	if err != nil {
		return 0, err
	}

	return newMiniURLID, nil
}

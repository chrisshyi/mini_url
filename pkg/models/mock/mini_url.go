package mock

import (
	"chrisshyi.net/mini_url/pkg/models"
)

// MiniURLModel is used for database operations
type MiniURLModel struct{}

var mockMiniURL = &models.MiniURL{
	ID:     1,
	URL:    "http://mock.com",
	Visits: 1,
}

// GetByID retrieves a MiniURL by its ID field
func (m *MiniURLModel) GetByID(ID int) (*models.MiniURL, error) {
	switch ID {
	case 1:
		return mockMiniURL, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// GetByURL retrieves a MiniURL by its url field
func (m *MiniURLModel) GetByURL(URL string) (*models.MiniURL, error) {
	switch URL {
	case "http://mock.com":
		return mockMiniURL, nil
	default:
		return nil, models.ErrNoRecord
	}
}

// Insert : insert a new URL
func (m *MiniURLModel) Insert(URL string) (int, error) {
	switch URL {
	case "http://mock.com":
		return 1, nil
	default:
		return 2, nil
	}
}

package postgres

import (
	"database/sql"
	"errors"

	"chrisshyi.net/snippetbox/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

// UserModel is used for interacting with the users table
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	emailExists, err := m.emailExists(email)
	if err != nil {
		return err
	}
	if emailExists {
		return models.ErrDuplicateEmail
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES ($1, $2, $3, NOW())`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))

	if err != nil {
		return err
	}
	return nil
}

// Authenticate verifies whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

// Get fetches details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	u := &models.User{}
	stmt := `SELECT id, name, email, created, active FROM users WHERE id = $1`
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created, &u.Active)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return u, nil
}

// emailExists checks whether a user with the given email exists
func (m *UserModel) emailExists(email string) (bool, error) {
	stmt := `SELECT email FROM users WHERE email = $1`
	var selectedEmail string
	err := m.DB.QueryRow(stmt, email).Scan(&selectedEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

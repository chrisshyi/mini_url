package models

import "errors"

var (
	// ErrNoRecord represents an error where a matching miniURL wasn't found
	ErrNoRecord = errors.New("models: no matching MiniURL found")
)

// MiniURL represents a shortened URL
type MiniURL struct {
	ID     int
	URL    string
	Visits int
}

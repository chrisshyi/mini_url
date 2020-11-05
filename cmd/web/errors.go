package main

import (
	"errors"
)

var (
	// ErrInvalidShortURL represents an invalid short URL error
	ErrInvalidShortURL = errors.New("Invalid short URL")
	// ErrInvalidID represents an invalid ID error
	ErrInvalidID = errors.New("Invalid ID")
)

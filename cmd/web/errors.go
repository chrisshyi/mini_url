package main

import (
	"errors"
)

var (
	ErrInvalidShortURL = errors.New("Invalid short URL")
	ErrInvalidID       = errors.New("Invalid ID")
)

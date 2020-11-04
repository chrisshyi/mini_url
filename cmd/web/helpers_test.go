package main

import (
	"errors"
	"testing"
)

func TestShortURLToID(t *testing.T) {
	tests := []struct {
		shortURL string
		wantID   int
		wantErr  error
	}{
		{"", -1, errors.New("input cannot be empty")},
		{"a", 1, nil},
		{"aa", 63, nil},
		{"你qwe好", -1, errors.New("input cannot contain non-alphanumerics")},
		{"0a", 3287, nil},
	}

	for _, test := range tests {
		ID, err := shortURLToID(test.shortURL)
		if test.wantErr != nil {
			if err == nil {
				t.Errorf("Expected error %s", test.wantErr)
			}
		} else {
			if ID != test.wantID {
				t.Errorf("Expected ID %d but got %d", test.wantID, ID)
			}
		}
	}
}

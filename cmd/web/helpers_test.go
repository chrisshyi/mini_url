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
		{"a", -1, errors.New("cannot contain all a's")},
		{"aa", -1, errors.New("cannot contain all a's")},
		{"aaaaaaaa", -1, errors.New("cannot contain all a's")},
		{"你qwe好", -1, errors.New("input cannot contain non-alphanumerics")},
		{"0a", 3224, nil},
		{"bb", 63, nil},
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

func TestIDToShortURL(t *testing.T) {
	tests := []struct {
		ID           int
		wantShortURL string
		wantErr      error
	}{
		{-1, "", errors.New("input cannot be less than 1")},
		{0, "", errors.New("input cannot be less than 1")},
		{3224, "0a", nil},
		{63, "bb", nil},
		{63210, "qBG", nil},
		{1, "b", nil},
		{3843, "99", nil},
	}

	for _, test := range tests {
		shortURL, err := IDToShortURL(test.ID)
		if test.wantErr != nil {
			if err == nil {
				t.Errorf("Expected error %s", test.wantErr)
			}
		} else {
			if shortURL != test.wantShortURL {
				t.Errorf("Expected shortURL %s but got %s", test.wantShortURL, shortURL)
			}
		}
	}

}

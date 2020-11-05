package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
)

// shortURLToID converts a short URL to a numerical ID
func shortURLToID(shortURL string) (int, error) {
	fmt.Printf("Input string: %s\n", shortURL)
	consecutiveAs, err := regexp.MatchString(`^a+$`, shortURL)
	if err != nil {
		return -1, err
	}
	if consecutiveAs {
		return -1, errors.New("shortURL cannot contain all a's")
	}
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, shortURL)
	fmt.Printf("%s matched: %v\n", shortURL, matched)
	if err != nil {
		return -1, err
	}
	if !matched {
		return -1, errors.New("shortURL cannot be empty")
	}
	multipliers := make([]int, 1)
	for i := 0; i < len(shortURL); i++ {
		ch := shortURL[i]
		if 'a' <= ch && ch <= 'z' {
			multipliers = append(multipliers, int(ch-'a'))
		} else if 'A' <= ch && ch <= 'Z' {
			multipliers = append(multipliers, int(ch-'A'+26))
		} else {
			multipliers = append(multipliers, int(ch-'0'+52))
		}
	}
	var ID int
	n := len(multipliers)
	for i := 0; i < n; i++ {
		ID += int(math.Pow(62, float64(n-i-1))) * multipliers[i]
	}
	return ID, nil
}

// IDToShortURL converts an ID to a short URL
func IDToShortURL(ID int) (string, error) {
	if ID < 1 {
		return "", errors.New("ID cannot be less than 1")
	}
	var digits []int
	num := ID
	for num > 0 {
		rem := num % 62
		digits = append(digits, rem)
		num = num / 62
	}
	// TODO: finish this!
	return "", nil
}

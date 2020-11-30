package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/golang/gddo/httputil/header"
)

// shortURLToID converts a short URL to a numerical ID
func shortURLToID(shortURL string) (int, error) {
	consecutiveAs, err := regexp.MatchString(`^a+$`, shortURL)
	if err != nil {
		return -1, err
	}
	if consecutiveAs {
		return -1, ErrInvalidShortURL
	}
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, shortURL)
	if err != nil {
		return -1, err
	}
	if !matched {
		return -1, ErrInvalidShortURL
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
		return "", errors.New("Invalid ID")
	}
	var digits []int
	num := ID
	for num > 0 {
		rem := num % 62
		digits = append(digits, rem)
		num = num / 62
	}
	var sb strings.Builder
	for i := len(digits) - 1; i > -1; i-- {
		digit := digits[i]
		if digit < 26 {
			sb.WriteByte('a' + byte(digit))
		} else if digit < 52 {
			sb.WriteByte('A' + byte(digit-26))
		} else {
			sb.WriteByte('0' + byte(digit-52))
		}
	}
	return sb.String(), nil
}

func (app *application) logErr(errorMsg string) {
	trace := fmt.Sprintf("%s\n%s", errorMsg, debug.Stack())
	app.errorLog.Output(2, trace)
}

func (app *application) logInfo(infoMsg string) {
	app.infoLog.Output(2, infoMsg)
}

func hasHTTPPrefix(URL string) (bool, error) {
	matched, err := regexp.MatchString(`^https?://.*`, URL)
	if err != nil {
		return false, err
	}
	return matched, nil
}

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

package utils

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateFullName(s string) error {
	var (
		totalChar = len(s)
	)

	if totalChar < 3 || totalChar > 60 {
		return errors.New("must be at minimum 3 characters and maximum 60 characters")
	}

	return nil
}

func ValidatePhoneNumbers(s string) error {
	var (
		totalChar = len(s)
		err       = make([]string, 0)
	)

	if totalChar < 10 || totalChar > 13 {
		err = append(err, "must be minimum 10 characters and maximum 13 characters")
	}

	if totalChar > 2 && !(s[0] == '+' && s[1] == '6' && s[2] == '2') {
		err = append(err, "must start with the Indonesia country code “+62”")
	}

	if len(err) == 0 {
		return nil
	}

	return errors.New(strings.Join(err, " & "))
}

func ValidatePassword(s string) error {
	var (
		totalChar                int
		capital, number, special bool
		err                      = make([]string, 0)
	)
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			capital = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		default:
		}
		totalChar++
	}

	if totalChar < 6 || totalChar > 64 {
		err = append(err, "must be minimum 6 characters and maximum 64 characters")
	}

	if !(number && capital && special) {
		err = append(err, "must containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")
	}

	if len(err) == 0 {
		return nil
	}

	return errors.New(strings.Join(err, " & "))
}

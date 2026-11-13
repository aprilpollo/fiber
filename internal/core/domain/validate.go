package domain

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

const maxTextLength = 255

func requireText(field, value string) error {
	switch {
	case strings.TrimSpace(value) == "":
		return invalid(field, "is required")
	case utf8.RuneCountInString(value) > maxTextLength:
		return invalid(field, "must be at most 255 characters")
	}
	return nil
}

func requireEmail(field, value string) error {
	if err := requireText(field, value); err != nil {
		return err
	}
	// ParseAddress also accepts `Name <a@b.c>`; only a bare address is an email here.
	if addr, err := mail.ParseAddress(value); err != nil || addr.Address != value {
		return invalid(field, "must be a valid email address")
	}
	return nil
}

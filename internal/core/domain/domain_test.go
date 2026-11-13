package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	u, err := NewUser("  Alice  ", " alice@example.com ")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Fatal("NewUser must assign an identity")
	}
	if u.Name != "Alice" || u.Email != "alice@example.com" {
		t.Fatalf("input not trimmed: %+v", u)
	}
}

func TestUserInvariants(t *testing.T) {
	tests := []struct {
		name, userName, email, field string
	}{
		{"blank name", "  ", "a@b.co", "name"},
		{"name too long", strings.Repeat("ก", 256), "a@b.co", "name"},
		{"blank email", "Alice", "", "email"},
		{"not an email", "Alice", "alice", "email"},
		{"display-name form", "Alice", "Alice <a@b.co>", "email"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewUser(tc.userName, tc.email)
			var verr *ValidationError
			if !errors.As(err, &verr) || verr.Field != tc.field {
				t.Fatalf("want ValidationError on %s, got %v", tc.field, err)
			}
		})
	}
}

func TestNameLengthCountsCharactersNotBytes(t *testing.T) {
	// 255 Thai characters are 765 bytes — still a valid name.
	if _, err := NewUser(strings.Repeat("ก", 255), "a@b.co"); err != nil {
		t.Fatal(err)
	}
}

func TestBookInvariants(t *testing.T) {
	if _, err := NewBook("Clean Architecture", "Robert C. Martin"); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ title, author, field string }{
		{"", "x", "title"},
		{"x", " ", "author"},
	} {
		var verr *ValidationError
		if _, err := NewBook(tc.title, tc.author); !errors.As(err, &verr) || verr.Field != tc.field {
			t.Fatalf("want ValidationError on %s, got %v", tc.field, err)
		}
	}
}

package lid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

const separator = "-"

type Prefix string

const (
	Bot Prefix = "bot"
)

var (
	ErrInvalidID     = errors.New("invalid identifier")
	ErrInvalidPrefix = errors.New("invalid identifier prefix")
)

func NewBot() string { return newID(Bot) }

func newID(prefix Prefix) string {
	return string(prefix) + separator + ulid.Make().String()
}

// Validate checks both the resource prefix and canonical ULID representation.
func Validate(value string, expected Prefix) error {
	_, err := Parse(value, expected)
	return err
}

// Parse validates the resource prefix and returns the underlying ULID.
func Parse(value string, expected Prefix) (ulid.ULID, error) {
	prefix, rawULID, found := strings.Cut(value, separator)
	if !found || Prefix(prefix) != expected {
		return ulid.ULID{}, fmt.Errorf("%w: expected %s prefix", ErrInvalidPrefix, expected)
	}
	if strings.Contains(rawULID, separator) {
		return ulid.ULID{}, fmt.Errorf("%w: unexpected separator", ErrInvalidID)
	}
	parsed, err := ulid.ParseStrict(rawULID)
	if err != nil {
		return ulid.ULID{}, fmt.Errorf("%w: %v", ErrInvalidID, err)
	}
	return parsed, nil
}

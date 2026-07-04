package utils

import (
	"fmt"

	"github.com/gofrs/uuid"
)

func NewUUID() (string, error) {
	id, err := uuid.NewV7() // Changed to NewV7. It is more in-line with the db I want to use. It is also more future-proof and has better sorting properties.
	if err != nil {
		return "", fmt.Errorf("could not generate UUID: %w", err)
	}
	return id.String(), nil // Converts the internal, raw 16-byte array representation of the UUID into the human-readable 36-character hyphenated string format aka xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
}
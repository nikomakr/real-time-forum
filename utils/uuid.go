package utils

import (
	"log"

	"github.com/gofrs/uuid"
)

func NewUUID() string {
	id, err := uuid.NewV4() // I am not sure about NewV4 Vs NewV7. I might have to revisit this later.
	if err != nil {
		log.Fatalf("could not generate UUID: %v", err)
	}
	return id.String() // Converts the internal, raw 16-byte array representation of the UUID into the human-readable 36-character hyphenated string format aka xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`.
}
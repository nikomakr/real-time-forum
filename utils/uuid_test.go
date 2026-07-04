package utils_test

import (
	"strings"
	"testing"

	"real-time-forum/utils"
)

func TestNewUUID_IsUnique(t *testing.T) { // The t holds the memory address of the testing.T, right? As we have a pointer! The T is a struct of the library!
	a := utils.NewUUID()
	b := utils.NewUUID()

	if a == b {
		t.Error("two UUIDs should not be identical")
	}
}

func TestNewUUID_IsValidFormat(t *testing.T) {
	id := utils.NewUUID()

	if len(id) != 36 {
		t.Errorf("expected 36 characters, got %d", len(id))
	}

	parts := strings.Split(id, "-")
	if len(parts) != 5 {
		t.Errorf("expected 5 hyphen-separated groups, got %d", len(parts))
	}
}
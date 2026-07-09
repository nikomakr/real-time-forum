package utils_test

import (
	"strings"
	"testing"

	"real-time-forum/utils"
)

func TestNewUUID_IsUnique(t *testing.T) {
	a, err := utils.NewUUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := utils.NewUUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == b {
		t.Error("two UUIDs should not be identical")
	}
}

func TestNewUUID_IsValidFormat(t *testing.T) {
	id, err := utils.NewUUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(id) != 36 {
		t.Errorf("expected 36 characters, got %d", len(id))
	}
	parts := strings.Split(id, "-")
	if len(parts) != 5 {
		t.Errorf("expected 5 hyphen-separated groups, got %d", len(parts))
	}
}

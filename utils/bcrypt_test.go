package utils_test

import (
	"testing"

	"real-time-forum/utils"
)

func TestHashPassword_ProducesHash(t *testing.T) {
	hash, err := utils.HashPassword("secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hash) == 0 {
		t.Error("expected a non-empty hash")
	}
}

func TestHashPassword_NeverStoresPlainText(t *testing.T) {
	password := "secret123"
	hash, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == password {
		t.Error("hash must not equal the original password")
	}
}

func TestHashPassword_TwoHashesOfSamePasswordDiffer(t *testing.T) {
	a, _ := utils.HashPassword("secret123")
	b, _ := utils.HashPassword("secret123")
	if a == b {
		t.Error("bcrypt must produce a different hash each time due to random salt")
	}
}

func TestCheckPassword_CorrectPasswordPasses(t *testing.T) {
	hash, _ := utils.HashPassword("secret123")
	err := utils.CheckPassword("secret123", hash)
	if err != nil {
		t.Errorf("expected correct password to pass, got: %v", err)
	}
}

func TestCheckPassword_WrongPasswordFails(t *testing.T) {
	hash, _ := utils.HashPassword("secret123")
	err := utils.CheckPassword("wrongpassword", hash)
	if err == nil {
		t.Error("expected wrong password to fail but it passed")
	}
}

func TestCheckPassword_EmptyPasswordFails(t *testing.T) {
	hash, _ := utils.HashPassword("secret123")
	err := utils.CheckPassword("", hash)
	if err == nil {
		t.Error("expected empty password to fail but it passed")
	}
}

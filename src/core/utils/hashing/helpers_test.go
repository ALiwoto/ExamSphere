package hashing_test

import (
	"testing"

	"OnlineExams/src/core/utils/hashing"
)

func TestSHA256(t *testing.T) {
	h := hashing.HashSHA256("hello")
	if h != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Error("Expected hash, got", h)
	}
}

func TestGenerateAccessHash(t *testing.T) {
	h := hashing.GenerateAuthHash()
	if len(h) != hashing.AuthHashSize {
		t.Error("Expected 8 characters, got", len(h))
	}
}

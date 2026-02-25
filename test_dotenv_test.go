package main

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDotenvPasswordComparison(t *testing.T) {
	// create a password and hash it, then expose via env to simulate dotenv
	pass := "s3cr3t-pass"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed generating hash: %v", err)
	}

	os.Setenv("ADMIN_PASS_HASH", string(hash))
	os.Setenv("MASTER_PASSWORD", pass)
	os.Setenv("MINT_PASS_HASH", "unused-mint-hash")

	// retrieve and compare
	gotHash := os.Getenv("ADMIN_PASS_HASH")
	gotPass := os.Getenv("MASTER_PASSWORD")

	if err := bcrypt.CompareHashAndPassword([]byte(gotHash), []byte(gotPass)); err != nil {
		t.Fatalf("expected password to match hash, but CompareHashAndPassword failed: %v", err)
	}

	// negative check: wrong password should not match
	os.Setenv("MASTER_PASSWORD", "wrong-pass")
	if err := bcrypt.CompareHashAndPassword([]byte(gotHash), []byte(os.Getenv("MASTER_PASSWORD"))); err == nil {
		t.Fatalf("expected mismatch for wrong password, but CompareHashAndPassword returned nil")
	}
}

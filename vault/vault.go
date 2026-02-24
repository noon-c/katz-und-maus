package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// This is a simple vault implementation that encrypts a text using a master password and stores it in an S3-compatible object storage (Cloudflare R2 in this case).
//
//go run ./vault write
//go run ./vault read

func main() {

	_ = godotenv.Load()

	if len(os.Args) < 2 {
		log.Fatal("usage: go run vault/vault.go [read|write]")
	}
	switch os.Args[1] {
	case "read":
		runRead()
	case "write":
		runWrite()
	default:
		log.Fatal("unknown command: " + os.Args[1])
	}
}

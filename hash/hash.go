package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// go get -u golang.org/x/crypto/bcrypt
// use: go run ./hash
func main() {
	// pass := "A8914875"
	pass := "123456"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	s := string(hash)
	s = strings.ReplaceAll(s, "$", "\\$")
	fmt.Println("\n=== hash \n", string(hash))
	fmt.Println("\n=== quoted hash \n\"" + s + "\" \n")
}

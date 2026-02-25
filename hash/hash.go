package main

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var pass = "123456"

// go get -u golang.org/x/crypto/bcrypt
// use: go run ./hash
func main() {
	// pass = "8914875"
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	s := string(hash)
	s = strings.ReplaceAll(s, "$", "\\$")
	fmt.Println("\n=== pass: ", string(pass))
	fmt.Println("\n=== hash: ", string(hash))
	fmt.Println("\n=== quoted hash: \"" + s + "\" \n")
}

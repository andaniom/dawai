package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// From DB: $2a$12$9Imjwp7N3Pg9JYrSFXGhBeS/CigEfaTPDumYSInWxuvfAshWNBcMq
	hash := "$2a$12$9Imjwp7N3Pg9JYrSFXGhBeS/CigEfaTPDumYSInWxuvfAshWNBcMq"
	password := "password"

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
	} else {
		fmt.Printf("OK: password matches hash\n")
	}
}

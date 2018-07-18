package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Printf("Input your password: ")
	var password []byte
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		password = scanner.Bytes()
		break
	}
	if len(password) == 0 {
		log.Fatal("Password is empty")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your hashed password: %s\n", string(hash))

}

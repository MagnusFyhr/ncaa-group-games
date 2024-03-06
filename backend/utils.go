// utils.go
package main

import (
	"fmt"
	"math/rand"
)

func generateRandomID() string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

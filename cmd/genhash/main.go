package main

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Dùng để tạo hash cho password mới.
// Usage: go run cmd/genhash/main.go
// Copy ADMIN_PASSWORD_HASH value vào .env
func main() {
	password := "admin123" // đổi password tại đây
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	b64 := base64.StdEncoding.EncodeToString(hash)
	fmt.Printf("ADMIN_PASSWORD_HASH=%s\n", b64)
}

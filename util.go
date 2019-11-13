package main

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a salted hash value of the specified password.
func HashPassword(password string) (hash string) {
	hashVal, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		panic(err)
	}

	return string(hashVal)
}

// ComparePasswordHashes compares an unhashed password with a hashed password.
func ComparePasswordHashes(password string, hashedPassword string) (equal bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return false
	}

	return true
}

// Converts little endian bytes to an int64.
func bytesToInt64(b []byte) int64 {
	var val int64

	val |= int64(b[0])
	val |= int64(b[1]) << 8
	val |= int64(b[2]) << 16
	val |= int64(b[3]) << 24

	return val
}

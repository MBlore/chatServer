package main

import (
	"bytes"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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

// GenerateGUID returns a new random GUID in the string format of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func GenerateGUID() string {
	u1 := uuid.Must(uuid.NewV4())
	return u1.String()
}

func readInt32(r *bytes.Reader) int64 {
	numBytes := make([]byte, 4)
	r.Read(numBytes)
	num := bytesToInt64(numBytes)

	return num
}

func readLenString(r *bytes.Reader) *string {
	len := readInt32(r)
	if len == 0 {
		return nil
	}

	strBytes := make([]byte, len)
	r.Read(strBytes)
	str := string(strBytes)

	return &str
}

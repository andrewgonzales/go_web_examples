package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

var time = uint32(3)
var memory = uint32(64 * 1024)
var threads = uint8(4)
var keyLen = uint32(32)

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func generateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPasswordArgon2(password string) (string, error) {
	salt, err := generateSalt(16)
	if err != nil {
		fmt.Println("Error generating salt: ", err)
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Base64 encode the hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, memory, time, threads, b64Salt, b64Hash)

	return encodedHash, nil
}

func decodeHash(encodedHash string) (salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		fmt.Println("Invalid hash")
		return nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		fmt.Println("Error scanning version: ", err)
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, ErrIncompatibleVersion
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		fmt.Println("Error decoding salt: ", err)
		return nil, nil, err
	}

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		fmt.Println("Error decoding hash: ", err)
		return nil, nil, err
	}

	return salt, hash, nil
}

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
func ComparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the salt and derived key from the encoded password
	// hash.
	salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "secret"

	// bcrypt

	hash, err := HashPassword(password)

	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	fmt.Println("bcrypt hash:    ", hash)

	match := CheckPasswordHash(password, hash)
	fmt.Println("bcrypt match:   ", match)

	// argon2
	hash, err = HashPasswordArgon2(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	fmt.Println("argon2 hash: ", hash)

	match, err = ComparePasswordAndHash(password, hash)
	if err != nil {
		fmt.Println("Error verifying password (argon2)", err)
		return
	}

	fmt.Printf("argon2 match: %v\n", match)

}

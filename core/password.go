package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	gotils "github.com/savsgio/gotils/strconv"
	"golang.org/x/crypto/argon2"
)

const (
	memory      = 19456
	iterations  = 2
	parallelism = 1
	saltLength  = 16
	keyLength   = 32
)

func PasswordHash(password string) string {
	salt := make([]byte, saltLength)
	rand.Read(salt)
	hash := argon2.IDKey(gotils.S2B(password), salt, iterations, memory, parallelism, keyLength)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash))
}

func VerifyPassword(password string, encodedHash string) bool {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		return false
	}

	var version int
	_, err := fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return false
	}
	if version != argon2.Version {
		return false
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(values[4])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(values[5])
	if err != nil {
		return false
	}
	keyLength := uint32(len(hash))

	passwordHash := argon2.IDKey(gotils.S2B(password), salt, iterations, memory, parallelism, keyLength)

	return subtle.ConstantTimeCompare(hash, passwordHash) == 1
}

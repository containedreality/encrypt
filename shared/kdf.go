package shared

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

func DeriveKeysPassword(password string, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 16)
		rand.Read(salt)
	}

	ek := argon2.IDKey([]byte(password), salt, 3, 256*1024, 4, 32)

	return ek, salt
}

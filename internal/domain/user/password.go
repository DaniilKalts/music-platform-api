package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const saltSize = 16
const minPasswordLength = 8

type Password struct {
	Hash string
	Salt string
}

func NewPassword(plain string) (Password, error) {
	if len(strings.TrimSpace(plain)) < minPasswordLength {
		return Password{}, ErrInvalidPassword
	}

	salt, err := generateSalt()
	if err != nil {
		return Password{}, fmt.Errorf("generate salt: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain+salt), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, fmt.Errorf("hash password: %w", err)
	}

	return Password{Hash: string(hash), Salt: salt}, nil
}

func (p Password) Matches(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plain+p.Salt)) == nil
}

func generateSalt() (string, error) {
	buf := make([]byte, saltSize)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

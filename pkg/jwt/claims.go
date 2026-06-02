package jwt

import (
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role,omitempty"`
	Type   TokenType `json:"type"`
	gojwt.RegisteredClaims
}

func newClaims(userID uuid.UUID, role string, tokenType TokenType, issuedAt, expiresAt time.Time) Claims {
	return Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
			IssuedAt:  gojwt.NewNumericDate(issuedAt),
			Issuer:    issuer,
			Subject:   userID.String(),
		},
	}
}

package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	gojwt "github.com/golang-jwt/jwt/v5"
)

const issuer = "music-platform-api"

var ErrInvalidToken = errors.New("invalid or expired token")

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role,omitempty"`
	Type   TokenType `json:"type"`
	gojwt.RegisteredClaims
}

type Pair struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewManager(accessSecret, refreshSecret []byte, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (m *Manager) GeneratePair(userID uuid.UUID, role string) (*Pair, error) {
	access, accessExpiresAt, err := m.generate(userID, role, TokenTypeAccess, m.accessSecret, m.accessTTL)
	if err != nil {
		return nil, err
	}

	refresh, refreshExpiresAt, err := m.generate(userID, "", TokenTypeRefresh, m.refreshSecret, m.refreshTTL)
	if err != nil {
		return nil, err
	}

	return &Pair{
		AccessToken:           access,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refresh,
		RefreshTokenExpiresAt: refreshExpiresAt,
	}, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (m *Manager) generate(userID uuid.UUID, role string, tokenType TokenType, secret []byte, ttl time.Duration) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(expiresAt),
			IssuedAt:  gojwt.NewNumericDate(now),
			Issuer:    issuer,
			Subject:   userID.String(),
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}

	return signed, expiresAt, nil
}

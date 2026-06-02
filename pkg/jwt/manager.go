package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	issuer = "music-platform-api"
)

type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewManager(accessSecret, refreshSecret []byte, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		accessSecret:  append([]byte(nil), accessSecret...),
		refreshSecret: append([]byte(nil), refreshSecret...),
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

func (m *Manager) ParseToken(tokenStr string, tokenType TokenType) (*Claims, error) {
	secret, ok := m.secretFor(tokenType)
	if !ok {
		return nil, ErrInvalidToken
	}

	token, err := gojwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		m.keyFunc(secret),
		gojwt.WithExpirationRequired(),
		gojwt.WithIssuer(issuer),
		gojwt.WithValidMethods([]string{gojwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != tokenType {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *Manager) generate(userID uuid.UUID, role string, tokenType TokenType, secret []byte, ttl time.Duration) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(ttl)
	claims := newClaims(userID, role, tokenType, now, expiresAt)

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}

	return signed, expiresAt, nil
}

func (m *Manager) secretFor(tokenType TokenType) ([]byte, bool) {
	switch tokenType {
	case TokenTypeAccess:
		return m.accessSecret, true
	case TokenTypeRefresh:
		return m.refreshSecret, true
	default:
		return nil, false
	}
}

func (m *Manager) keyFunc(secret []byte) gojwt.Keyfunc {
	return func(token *gojwt.Token) (interface{}, error) {
		if token.Method != gojwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	}
}

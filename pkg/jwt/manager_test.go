package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	gojwt "github.com/golang-jwt/jwt/v5"
)

func TestManagerParseToken(t *testing.T) {
	t.Parallel()

	manager := NewManager([]byte("access-secret"), []byte("refresh-secret"), time.Minute, time.Hour)
	userID := uuid.New()

	pair, err := manager.GeneratePair(userID, "USER")
	require.NoError(t, err)

	claims, err := manager.ParseToken(pair.AccessToken, TokenTypeAccess)
	require.NoError(t, err)
	require.Equal(t, userID, claims.UserID)
	require.Equal(t, "USER", claims.Role)
	require.Equal(t, TokenTypeAccess, claims.Type)

	_, err = manager.ParseToken(pair.AccessToken, TokenTypeRefresh)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestManagerParseTokenRejectsUnexpectedAlgorithm(t *testing.T) {
	t.Parallel()

	manager := NewManager([]byte("access-secret"), []byte("refresh-secret"), time.Minute, time.Hour)
	claims := Claims{
		UserID: uuid.New(),
		Role:   "USER",
		Type:   TokenTypeAccess,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Minute)),
			Issuer:    issuer,
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS384, claims)
	signed, err := token.SignedString(manager.accessSecret)
	require.NoError(t, err)

	_, err = manager.ParseToken(signed, TokenTypeAccess)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestManagerParseTokenRejectsMissingExpiration(t *testing.T) {
	t.Parallel()

	manager := NewManager([]byte("access-secret"), []byte("refresh-secret"), time.Minute, time.Hour)
	claims := Claims{
		UserID: uuid.New(),
		Role:   "USER",
		Type:   TokenTypeAccess,
		RegisteredClaims: gojwt.RegisteredClaims{
			Issuer: issuer,
		},
	}
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(manager.accessSecret)
	require.NoError(t, err)

	_, err = manager.ParseToken(signed, TokenTypeAccess)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestNewManagerCopiesSecrets(t *testing.T) {
	t.Parallel()

	accessSecret := []byte("access-secret")
	refreshSecret := []byte("refresh-secret")
	manager := NewManager(accessSecret, refreshSecret, time.Minute, time.Hour)
	accessSecret[0] = 'x'
	refreshSecret[0] = 'x'

	pair, err := manager.GeneratePair(uuid.New(), "USER")
	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)
	require.NotEmpty(t, pair.RefreshToken)
}

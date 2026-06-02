package httpx

import (
	"context"

	"github.com/google/uuid"
)

type requestIDKey struct{}
type userIdentityKey struct{}

type UserIdentity struct {
	ID   uuid.UUID
	Role string
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

func WithUser(ctx context.Context, identity UserIdentity) context.Context {
	return context.WithValue(ctx, userIdentityKey{}, identity)
}

func UserFromContext(ctx context.Context) (UserIdentity, bool) {
	identity, ok := ctx.Value(userIdentityKey{}).(UserIdentity)
	return identity, ok
}

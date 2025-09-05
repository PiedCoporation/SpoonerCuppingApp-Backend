package abstractions

import (
	"context"
)

type UserAuthUow interface {
	Begin(ctx context.Context) (UserAuthRepoProvider, error)
	Commit() error
	Rollback() error
}

type UserAuthRepoProvider interface {
	UserRepository() UserRepository
	RefreshTokenRepository() RefreshTokenRepository
}

package uow

import (
	"backend/internal/usecase/repository"
	"context"
)

type UserAuthUow interface {
	Begin(ctx context.Context) (UserAuthRepoProvider, error)
	Commit() error
	Rollback() error
}

type UserAuthRepoProvider interface {
	UserRepository() repository.UserRepository
	RefreshTokenRepository() repository.RefreshTokenRepository
}

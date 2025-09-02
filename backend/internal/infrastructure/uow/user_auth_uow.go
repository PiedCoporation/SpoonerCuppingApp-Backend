package uow

import (
	"backend/internal/infrastructure/repository/postgres"
	"backend/internal/usecase/repository"
	"backend/internal/usecase/uow"
	"context"

	"gorm.io/gorm"
)

type userAuthUow struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewUserAuthUow(db *gorm.DB) uow.UserAuthUow {
	return &userAuthUow{db: db}
}

// Begin implements uow.UserAuthUow.
func (u *userAuthUow) Begin(ctx context.Context) (uow.UserAuthRepoProvider, error) {
	tx := u.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	u.tx = tx
	return &userAuthRepoProvider{tx: tx}, nil
}

// Commit implements uow.UserAuthUow.
func (u *userAuthUow) Commit() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Commit().Error
}

// Rollback implements uow.UserAuthUow.
func (u *userAuthUow) Rollback() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Rollback().Error
}

// ===== RepoProvider =====
type userAuthRepoProvider struct {
	tx               *gorm.DB
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func (r *userAuthRepoProvider) UserRepository() repository.UserRepository {
	if r.userRepo == nil {
		r.userRepo = postgres.NewUserRepo(r.tx)
	}
	return r.userRepo
}

func (r *userAuthRepoProvider) RefreshTokenRepository() repository.RefreshTokenRepository {
	if r.refreshTokenRepo == nil {
		r.refreshTokenRepo = postgres.NewRefreshTokenRepo(r.tx)
	}
	return r.refreshTokenRepo
}

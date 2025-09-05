package postgres

import (
	"backend/internal/persistents/abstractions"
	"context"

	"gorm.io/gorm"
)

type userAuthUow struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewUserAuthUow(db *gorm.DB) abstractions.UserAuthUow {
	return &userAuthUow{db: db}
}

// Begin implements uow.UserAuthUow.
func (u *userAuthUow) Begin(ctx context.Context) (abstractions.UserAuthRepoProvider, error) {
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
	userRepo         abstractions.UserRepository
	refreshTokenRepo abstractions.RefreshTokenRepository
}

func (r *userAuthRepoProvider) UserRepository() abstractions.UserRepository {
	if r.userRepo == nil {
		r.userRepo = NewUserRepo(r.tx)
	}
	return r.userRepo
}

func (r *userAuthRepoProvider) RefreshTokenRepository() abstractions.RefreshTokenRepository {
	if r.refreshTokenRepo == nil {
		r.refreshTokenRepo = NewRefreshTokenRepo(r.tx)
	}
	return r.refreshTokenRepo
}

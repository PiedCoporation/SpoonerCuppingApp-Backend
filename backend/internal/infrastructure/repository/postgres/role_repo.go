package postgres

import (
	"backend/internal/domain/entities"
	"backend/internal/usecase/repository"

	"gorm.io/gorm"
)

type rolePgRepo struct {
	*genericRepository[entities.Role]
}

func NewRoleRepo(db *gorm.DB) repository.RoleRepository {
	return &rolePgRepo{
		genericRepository: NewGenericRepository[entities.Role](db),
	}
}

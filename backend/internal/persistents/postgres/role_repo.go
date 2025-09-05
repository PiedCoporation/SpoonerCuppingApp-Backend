package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type rolePgRepo struct {
	*genericRepository[entities.Role]
}

func NewRoleRepo(db *gorm.DB) abstractions.RoleRepository {
	return &rolePgRepo{
		genericRepository: NewGenericRepository[entities.Role](db),
	}
}

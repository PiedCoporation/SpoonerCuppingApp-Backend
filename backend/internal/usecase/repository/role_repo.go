package repository

import (
	"backend/internal/domain/entities"
)

type RoleRepository interface {
	GenericRepository[entities.Role]
}

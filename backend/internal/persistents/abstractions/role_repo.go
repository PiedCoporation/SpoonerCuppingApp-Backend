package abstractions

import (
	"backend/internal/domains/entities"
)

type RoleRepository interface {
	GenericRepository[entities.Role]
}

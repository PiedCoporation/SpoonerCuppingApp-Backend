package abstractions

import (
	"backend/internal/domains/entities"
)

type IRoleRepository interface {
	IGenericRepository[entities.Role]
}

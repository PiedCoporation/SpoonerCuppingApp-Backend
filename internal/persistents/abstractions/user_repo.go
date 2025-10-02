package abstractions

import (
	"backend/internal/domains/entities"
	"context"

	"github.com/google/uuid"
)

type IUserRepository interface {
	IGenericRepository[entities.User]
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	IsPhoneTaken(ctx context.Context, phone string, excludeUserID uuid.UUID) (bool, error)
	IsEmailTaken(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
}

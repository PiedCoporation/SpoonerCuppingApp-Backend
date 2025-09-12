package abstractions

import (
	"context"

	"github.com/google/uuid"
)

type GenericRepository[T any] interface {
	GetAll(ctx context.Context, preloads ...string) ([]T, error)
	GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error)
	Create(ctx context.Context, entity *T) error
	CreateRange(ctx context.Context, entities []T) error
	Update(ctx context.Context, id uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

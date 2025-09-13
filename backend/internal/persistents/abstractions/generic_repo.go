package abstractions

import (
	"context"

	"github.com/google/uuid"
)

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse represents a paginated response

type GenericRepository[T any] interface {
	GetAll(ctx context.Context, preloads ...string) ([]T, error)
	GetByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error)
	GetSingle(ctx context.Context, query string, preloads ...string) (*T, error)
	Create(ctx context.Context, entity *T) error
	CreateRange(ctx context.Context, entities []T) error
	Update(ctx context.Context, id uuid.UUID, fields map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

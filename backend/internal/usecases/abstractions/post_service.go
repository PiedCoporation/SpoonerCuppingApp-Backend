package abstractions

import (
	"backend/internal/contracts/common"
	"backend/internal/contracts/post"
	"context"
)

type IPostService interface {
	// Create(ctx context.Context, req event.CreatePostReq) error
	GetAll(ctx context.Context, pageSize int, pageNumber int, searchTerm string) (*common.PageResult[post.PostResponse], error)
	// GetByID(ctx context.Context, id uuid.UUID) (*event.GetEventByIDResponse, error)
	// Update(ctx context.Context, id uuid.UUID, event *entities.Event) error
	// Delete(ctx context.Context, id uuid.UUID) error
}

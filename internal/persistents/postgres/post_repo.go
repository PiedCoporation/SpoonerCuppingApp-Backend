package postgres

import (
	"backend/internal/domains/entities"
	"backend/internal/persistents/abstractions"

	"gorm.io/gorm"
)

type postPgRepo struct {
	*genericRepository[entities.Post]
}

func NewPostRepo(db *gorm.DB) abstractions.IPostRepository {
	return &postPgRepo{
		genericRepository: NewGenericRepository[entities.Post](db),
	}
}

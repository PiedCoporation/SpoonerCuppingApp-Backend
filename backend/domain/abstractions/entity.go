package abstractions

import (
	"github.com/google/uuid"
)

type Entity struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsDeleted bool      `gorm:"default:false"`
}

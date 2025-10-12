package entities

import (
	"backend/internal/domains/commons"

	"github.com/google/uuid"
)

type UserSampleTasting struct {
	commons.Entity
	ParentName     string `gorm:"column:parent_name"`
	ChildName      string `gorm:"column:child_name"`
	GrandChildName string `gorm:"column:grand_child_name"`
	commons.Auditable

	UserSampleID uuid.UUID   `gorm:"not null;index;column:user_sample_id"`
	UserSample   UserSample `gorm:"foreignKey:UserSampleID"`
}

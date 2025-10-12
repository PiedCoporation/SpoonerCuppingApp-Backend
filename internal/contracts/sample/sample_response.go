package sample

import (
	"backend/internal/constants/enums/processing"
	"backend/internal/constants/enums/roastinglever"
	"time"

	"github.com/google/uuid"
)

type SampleRes struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name            string                          `json:"name" example:"Coffee Sample"`
	RoastingDate    time.Time                       `json:"roasting_date" example:"2024-01-15T10:00:00Z"`
	RoastLevel      roastinglever.RoastingLeverEnum `json:"roast_level" example:"Medium"`
	AltitudeGrow    string                          `json:"altitude_grow" example:"1000m"`
	RoasteryName    string                          `json:"roastery_name" example:"Roastery Name"`
	RoasteryAddress string                          `json:"roastery_address" example:"Roastery Address"`
	BreedName       string                          `json:"breed_name" example:"Breed Name"`
	PreProcessing   processing.ProcessingEnum       `json:"pre_processing" example:"Washed"`
	GrowNation      string                          `json:"grow_nation" example:"Vietnam"`
	GrowAddress     string                          `json:"grow_address" example:"Grow Address"`
	Price           float64                         `json:"price" example:"100000"`
	UserSampleTastings []UserSampleTastingRes `json:"user_sample_tastings" example:"User Sample Tastings"`
}

type UserSampleTastingRes struct {
	ParentName string `json:"parent_name" example:"Parent Name"`
	ChildName string `json:"child_name" example:"Child Name"`
	GrandChildName string `json:"grand_child_name" example:"Grand Child Name"`
}
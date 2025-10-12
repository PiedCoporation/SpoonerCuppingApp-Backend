package sample

import (
	"backend/internal/constants/enums/processing"
	"backend/internal/constants/enums/roastinglever"
	"time"
)

type SampleReq struct {
	Name string
	RoastingDate time.Time
	RoastLevel roastinglever.RoastingLeverEnum
	AltitudeGrow string
	RoasteryName string
	RoasteryAddress string
	BreedName string
	PreProcessing processing.ProcessingEnum
	GrowNation string
	GrowAddress string
	Price float64
	UserSampleTastings []UserSampleTastingReq
}


type UserSampleTastingReq struct {
	ParentName string
	ChildName string
	GrandChildName string
}


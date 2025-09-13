package event

import (
	"backend/internal/constants/enums/processing"
	"backend/internal/constants/enums/roastinglever"
	"time"
)

type CreateEventReq struct {
	Name string
	DateOfEvent time.Time
	StartTime time.Time
	EndTime time.Time
	Limit int
	NumberSamples int
	PhoneContact string
	EmailContact string
	IsPublic bool
	Samples []NewSampleReq
	EventAddress []NewEventAddressReq
}

type UpdateEventReq struct {
	Name string
	DateOfEvent time.Time
	StartTime time.Time
	EndTime time.Time
	Limit int
	NumberSamples int
	PhoneContact string
	EmailContact string
	IsPublic bool
}

type NewSampleReq struct {
	Name            string  
	RoastingDate    time.Time                      
	RoastLevel      roastinglever.RoastingLeverEnum 
	AltitudeGrow    string                          
	RoasteryName    string                          
	RoasteryAddress string                          
	BreedName       string                         
	PreProcessing   processing.ProcessingEnum      
	GrowNation      string                          
	GrowAddress     string                          
	Price           float64     
}

type NewEventAddressReq struct {
	Province string
	District string
	Longitude string
	Latitude string
	Ward string
	Street string
	Phone string
}
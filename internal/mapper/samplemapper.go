package mapper

import (
	sampleContract "backend/internal/contracts/sample"
	"backend/internal/domains/entities"
)

func MapUserSampleToContractGetAllSampleResponse(u *entities.UserSample) *sampleContract.SampleRes {
	// Map tasting notes
	tastings := make([]sampleContract.UserSampleTastingRes, len(u.UserSampleTastings))
	for i, tasting := range u.UserSampleTastings {
		tastings[i] = sampleContract.UserSampleTastingRes{
			ParentName:     tasting.ParentName,
			ChildName:      tasting.ChildName,
			GrandChildName: tasting.GrandChildName,
		}
	}

	return &sampleContract.SampleRes{
		ID:        u.ID,
		Name: u.Name,
		RoastingDate: u.RoastingDate,
		RoastLevel: u.RoastLevel,
		AltitudeGrow: u.AltitudeGrow,
		RoasteryName: u.RoasteryName,
		RoasteryAddress: u.RoasteryAddress,
		BreedName: u.BreedName,
		PreProcessing: u.PreProcessing,
		GrowNation: u.GrowNation,
		GrowAddress: u.GrowAddress,
		Price: u.Price,
		UserSampleTastings: tastings,
	}
}
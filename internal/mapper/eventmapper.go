package mapper

import (
	"backend/internal/constants/enums/eventregisterstatus"
	eventContract "backend/internal/contracts/event"
	"backend/internal/domains/entities"
)

func MapEventToContractGetAllEventResponse(e *entities.Event) eventContract.Event {
	var addresses []eventContract.EventAddress
	for _, addr := range e.EventAddress {
		addresses = append(addresses, eventContract.EventAddress{
			Province:  addr.Province,
			District:  addr.District,
			Longitude: addr.Longitude,
			Latitude:  addr.Latitude,
			Ward:      addr.Ward,
			Street:    addr.Street,
			Phone:     addr.Phone,
		})
	}

	return eventContract.Event{
		ID:             e.ID,
		Name:           e.Name,
		DateOfEvent:    e.DateOfEvent,
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		Limit:          e.Limit,
		TotalCurrent:   e.TotalCurrent,
		NumberSamples:  e.NumberSamples,
		PhoneContact:   e.PhoneContact,
		EmailContact:   e.EmailContact,
		IsPublic:       e.IsPublic,
		RegisterDate:   e.RegisterDate,
		RegisterStatus: eventregisterstatus.RegisterStatusEnum(e.RegisterStatus),
		EventAddress:   addresses,
		HostBy: eventContract.HostBy{
			ID: e.HostBy.ID,
			FirstName: e.HostBy.FirstName,
			LastName: e.HostBy.LastName,
			Email: e.HostBy.Email,
			Phone: e.HostBy.Phone,
		},
	}
}

func MapEventToContractGetEventByIDResponse(e *entities.Event) eventContract.GetEventByIDResponse {
	var addresses []eventContract.EventAddress
	for _, addr := range e.EventAddress {
		addresses = append(addresses, eventContract.EventAddress{
			Province:  addr.Province,
			District:  addr.District,
			Longitude: addr.Longitude,
			Latitude:  addr.Latitude,
			Ward:      addr.Ward,
			Street:    addr.Street,
			Phone:     addr.Phone,
		})
	}

	var samples []eventContract.EventSample
	for _, sample := range e.EventSamples {
		samples = append(samples, eventContract.EventSample{
			ID: sample.ID,
			Name: sample.UserSample.Name,
			RoastingDate: sample.UserSample.RoastingDate,
			RoastLevel: sample.UserSample.RoastLevel,
			AltitudeGrow: sample.UserSample.AltitudeGrow,
			RoasteryName: sample.UserSample.RoasteryName,
			RoasteryAddress: sample.UserSample.RoasteryAddress,
			BreedName: sample.UserSample.BreedName,
			PreProcessing: sample.UserSample.PreProcessing,
			GrowNation: sample.UserSample.GrowNation,
			GrowAddress: sample.UserSample.GrowAddress,
			Price: sample.Price,
			Rating: sample.Rating,
		})
	}

	return eventContract.GetEventByIDResponse{
		Event: eventContract.Event{
			ID:             e.ID,
			Name:           e.Name,
			DateOfEvent:    e.DateOfEvent,
			StartTime:      e.StartTime,
			EndTime:        e.EndTime,
			Limit:          e.Limit,
			TotalCurrent:   e.TotalCurrent,
			NumberSamples:  e.NumberSamples,
			PhoneContact:   e.PhoneContact,
			EmailContact:   e.EmailContact,
			IsPublic:       e.IsPublic,
			RegisterDate:   e.RegisterDate,
			RegisterStatus: eventregisterstatus.RegisterStatusEnum(e.RegisterStatus),
			EventAddress:   addresses,
			HostBy: eventContract.HostBy{
				ID: e.HostBy.ID,
				FirstName: e.HostBy.FirstName,
				LastName: e.HostBy.LastName,
				Email: e.HostBy.Email,
				Phone: e.HostBy.Phone,
			},
		},
		Samples: samples,
	}
}
//go:build wireinject

package event

import (
	"backend/internal/persistents/postgres"
	"backend/internal/presentations/http/v1/controller"
	"backend/internal/usecases"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitEventRouterHandler(
	db *gorm.DB,
) (*controller.EventController, error) {
	wire.Build(
		postgres.NewEventRepo,
		postgres.NewEventAddressRepo,
		postgres.NewEventSampleRepo,
		postgres.NewEventUserRepo,
		postgres.NewEventUow,
		usecases.NewEventService,
		controller.NewEventController,
	)
	return &controller.EventController{}, nil
}

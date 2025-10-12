//go:build wireinject

//go:generate wire

package sample

import (
	"backend/internal/persistents/postgres"
	"backend/internal/presentations/http/v1/controller"
	"backend/internal/usecases"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitSampleRouterHandler(
	db *gorm.DB,
) (*controller.SampleController, error) {
	wire.Build(
		postgres.NewSampleRepo,
		postgres.NewSampleUow,
		usecases.NewSampleService,
		controller.NewSampleController,
	)
	return &controller.SampleController{}, nil
}

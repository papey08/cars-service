package app

import (
	"cars-service/internal/adapters/api"
	"cars-service/internal/model"
	"cars-service/internal/repo"
	"cars-service/pkg/logger"
	"context"
)

type App interface {
	GetCarById(ctx context.Context, id uint64) (model.Car, error)
	GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error)

	AddCars(ctx context.Context, regNums []string) ([]model.Car, error)
	UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error)
	DeleteCar(ctx context.Context, id uint64) error
}

func New(r repo.Repo, cli api.Api, logs logger.Logger) App {
	return &appImpl{
		Logger: logs,
		Api:    cli,
		Repo:   r,
	}
}

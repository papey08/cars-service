package repo

import (
	"cars-service/internal/model"
	"context"
)

type Repo interface {
	GetCarById(ctx context.Context, id uint64) (model.Car, error)
	GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error)

	AddCar(ctx context.Context, car model.Car) (model.Car, error)
	UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error)
	DeleteCar(ctx context.Context, id uint64) error
}

package app

import (
	"cars-service/internal/model"
	"context"
)

type App interface {
	GetCarById(ctx context.Context, id uint64) (model.Car, error)
	GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error)

	AddCars(ctx context.Context, regNums []string) ([]model.Car, error)
	UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error)
	DeleteCar(ctx context.Context, id uint64) error
}

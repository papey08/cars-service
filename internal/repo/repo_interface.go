package repo

import (
	"cars-service/internal/model"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo is an interface of the database
type Repo interface {
	GetCarById(ctx context.Context, id uint64) (model.Car, error)
	GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error)

	AddCar(ctx context.Context, car model.Car) (model.Car, error)
	UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error)
	DeleteCar(ctx context.Context, id uint64) error
}

// New creates Repo implementation
func New(pool *pgxpool.Pool) Repo {
	return &repoImpl{
		Pool: pool,
	}
}

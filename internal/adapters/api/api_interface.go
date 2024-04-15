package api

import (
	"cars-service/internal/model"
	"context"
)

type Api interface {
	GetInfo(ctx context.Context, regNum string) (model.Car, error)
}

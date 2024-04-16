package api

import (
	"cars-service/internal/model"
	"context"
	"net/http"
)

// Api is an interface fo outer API
type Api interface {
	// GetInfo gets information about the car by its regNum
	GetInfo(ctx context.Context, regNum string) (model.Car, error)
}

// New creates Api implementation
func New(url string) Api {
	return &apiImpl{
		url:    url,
		Client: http.Client{},
	}
}

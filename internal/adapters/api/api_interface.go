package api

import (
	"cars-service/internal/model"
	"context"
	"net/http"
)

type Api interface {
	GetInfo(ctx context.Context, regNum string) (model.Car, error)
}

func NewApi(url string) Api {
	return &apiImpl{
		url:    url,
		Client: http.Client{},
	}
}

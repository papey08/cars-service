package app

import (
	"cars-service/internal/adapters/api"
	"cars-service/internal/model"
	"cars-service/internal/repo"
	"cars-service/pkg/logger"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"sync"
)

type appImpl struct {
	logger.Logger
	api.Api
	repo.Repo
}

func (a *appImpl) GetCarById(ctx context.Context, id string) (model.Car, error) {
	var err error
	defer a.writeLogs(logger.Fields{
		"Method": "GetCarById",
		"CarId":  id,
	}, err)

	car, err := a.Repo.GetCarById(ctx, id)
	return car, err
}

func (a *appImpl) GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error) {
	var err error
	defer a.writeLogs(logger.Fields{
		"Method": "GetCars",
	}, err)

	cars, err := a.Repo.GetCars(ctx, filter)
	return cars, err
}

func (a *appImpl) AddCars(ctx context.Context, regNums []string) ([]model.Car, error) {
	var err error
	defer a.writeLogs(logger.Fields{
		"Method": "AddCars",
	}, err)

	cars := make(map[string]model.Car, len(regNums))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, regNum := range regNums {
		wg.Add(1)
		go func(regNum string) {
			defer wg.Done()
			car, err := a.Api.GetInfo(ctx, regNum)
			if err != nil {
				a.Logger.Debug(logger.Fields{
					"Method": "GetInfo",
					"RegNum": regNum,
				}, err.Error())
				return
			}

			mu.Lock()
			cars[regNum] = car
			mu.Unlock()

			a.Logger.Debug(logger.Fields{
				"Method": "GetInfo",
				"RegNum": regNum,
			}, "ok")
		}(regNum)
	}
	wg.Wait()

	res := make([]model.Car, 0, len(cars))

	var gr errgroup.Group
	for _, car := range cars {
		car := car
		gr.Go(func() error {
			car, err := a.Repo.AddCar(ctx, car)
			if err != nil {
				return err
			}

			mu.Lock()
			res = append(res, car)
			mu.Unlock()
			return nil
		})
	}

	if err = gr.Wait(); err != nil {
		return []model.Car{}, err
	}
	return res, nil
}

func (a *appImpl) UpdateCar(ctx context.Context, id uint64, car model.Car) (model.Car, error) {
	var err error
	defer a.writeLogs(logger.Fields{
		"Method": "GetCars",
		"CarId":  id,
	}, err)

	car, err = a.Repo.UpdateCar(ctx, id, car)
	return car, err
}

func (a *appImpl) DeleteCar(ctx context.Context, id uint64) error {
	var err error
	defer a.writeLogs(logger.Fields{
		"Method": "DeleteCar",
		"CarId":  id,
	}, err)

	err = a.Repo.DeleteCar(ctx, id)
	return err
}

func (a *appImpl) writeLogs(fields logger.Fields, err error) {
	if errors.Is(err, model.ErrApiError) || errors.Is(err, model.ErrDatabaseError) {
		a.Logger.Error(fields, err.Error())
	} else if err != nil {
		a.Logger.Info(fields, err.Error())
	} else {
		a.Logger.Info(fields, "ok")
	}
}

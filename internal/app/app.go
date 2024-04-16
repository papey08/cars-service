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
	"time"
	"unicode"
)

type appImpl struct {
	logger.Logger
	api.Api
	repo.Repo
}

func (a *appImpl) GetCarById(ctx context.Context, id uint64) (model.Car, error) {
	var err error
	defer func() {
		a.writeLogs(logger.Fields{
			"Method": "GetCarById",
			"CarId":  id,
		}, err)
	}()

	car, err := a.Repo.GetCarById(ctx, id)
	return car, err
}

func (a *appImpl) GetCars(ctx context.Context, filter model.Filter) ([]model.Car, error) {
	var err error
	defer func() {
		a.writeLogs(logger.Fields{
			"Method": "GetCars",
		}, err)
	}()

	cars, err := a.Repo.GetCars(ctx, filter)
	return cars, err
}

func (a *appImpl) AddCars(ctx context.Context, regNums []string) ([]model.Car, error) {
	var err error
	defer func() {
		a.writeLogs(logger.Fields{
			"Method": "AddCars",
		}, err)
	}()

	cars := make(map[string]model.Car, len(regNums))
	var mu sync.Mutex
	var wg sync.WaitGroup

	// getting data of all regNums
	for _, regNum := range regNums {
		wg.Add(1)
		if !isValidRegNum([]rune(regNum)) {
			a.Logger.Debug(logger.Fields{
				"Method": "AddCars",
				"RegNum": regNum,
			}, model.ErrValidation.Error())
			continue
		}
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

	// adding all new cars to the database
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
	defer func() {
		a.writeLogs(logger.Fields{
			"Method": "UpdateCar",
			"CarId":  id,
		}, err)
	}()
	if !isValidRegNum([]rune(car.RegNum)) || !isValidYear(car.Year) {
		err = model.ErrValidation
		return model.Car{}, err
	}

	car, err = a.Repo.UpdateCar(ctx, id, car)
	return car, err
}

func (a *appImpl) DeleteCar(ctx context.Context, id uint64) error {
	var err error
	defer func() {
		a.writeLogs(logger.Fields{
			"Method": "DeleteCar",
			"CarId":  id,
		}, err)
	}()

	err = a.Repo.DeleteCar(ctx, id)
	return err
}

func isValidRegNum(regNum []rune) bool {
	if len(regNum) != 9 && len(regNum) != 8 {
		return false
	}
	if len(regNum) == 9 && !unicode.IsDigit(regNum[8]) {
		return false
	}

	return unicode.IsLetter(regNum[0]) &&
		unicode.IsDigit(regNum[1]) &&
		unicode.IsDigit(regNum[2]) &&
		unicode.IsDigit(regNum[3]) &&
		unicode.IsLetter(regNum[4]) &&
		unicode.IsLetter(regNum[5]) &&
		unicode.IsDigit(regNum[6]) &&
		unicode.IsDigit(regNum[7])
}

func isValidYear(year int) bool {
	return year >= 1900 && year <= time.Now().UTC().Year()
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

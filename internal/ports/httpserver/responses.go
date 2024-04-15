package httpserver

import "cars-service/internal/model"

func errorResponse(err error) carResponse {
	if err == nil {
		return carResponse{}
	}
	errStr := err.Error()
	return carResponse{
		Data: nil,
		Err:  &errStr,
	}
}

func carToCarData(car model.Car) carData {
	return carData{
		Id:     car.Id,
		RegNum: car.RegNum,
		Mark:   car.Mark,
		Model:  car.Model,
		Year:   car.Year,
		Owner: ownerData{
			Name:       car.Owner.Name,
			Surname:    car.Owner.Surname,
			Patronymic: car.Owner.Patronymic,
		},
	}
}

func carsToCarsData(cars []model.Car) []carData {
	data := make([]carData, len(cars))
	for i, car := range cars {
		data[i] = carToCarData(car)
	}
	return data
}

type carResponse struct {
	Data *carData `json:"data"`
	Err  *string  `json:"error"`
}

type carsResponse struct {
	Data []carData `json:"data"`
	Err  *string   `json:"error"`
}

type carData struct {
	Id     uint64    `json:"id"`
	RegNum string    `json:"regNum"`
	Mark   string    `json:"mark"`
	Model  string    `json:"model"`
	Year   int       `json:"year"`
	Owner  ownerData `json:"owner"`
}

type ownerData struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

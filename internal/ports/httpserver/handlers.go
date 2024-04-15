package httpserver

import (
	"cars-service/internal/app"
	"cars-service/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Получение информации об автомобиле по id
// @Description Возвращает информацию об автомобиле и его владельце, если автомобиль с указанным id существует
// @Produce json
// @Param id path int true "id автомобиля"
// @Success 200 {object} carResponse "Успешное получение информации"
// @Failure 400 {object} carResponse "Неверный формат входных данных"
// @Failure 404 {object} carResponse "Автомобиль с указанным id не найден"
// @Failure 500 {object} carResponse "Ошибка на стороне сервера"
// @Router /cars/{id} [get]
func handleGetCarById(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}

		car, err := a.GetCarById(c, id)

		switch {
		case err == nil:
			data := carToCarData(car)
			c.JSON(http.StatusOK, carResponse{
				Data: &data,
				Err:  nil,
			})
		case errors.Is(err, model.ErrCarNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse(model.ErrCarNotFound))
		case errors.Is(err, model.ErrDatabaseError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrDatabaseError))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrServiceError))
		}
	}
}

// @Summary Получение списка автомобилей
// @Description Возвращает список автомобилей, поддерживаются фильтрация и пагинация
// @Produce json
// @Param limit query int true ""
// @Param offset query int true ""
// @Param regNum query string false "Регистрационный номер"
// @Param mark query string false "Марка"
// @Param model query string false "Модель"
// @Param ownerName string false "Имя владельца"
// @Param ownerSurname string false "Фамилия владельца"
// @Param ownerPatronymic string false "Отчество владельца"
// @Success 200 {object} carResponse "Успешное получение информации"
// @Failure 400 {object} carResponse "Неверный формат входных данных"
// @Failure 500 {object} carResponse "Ошибка на стороне сервера"
// @Router /cars [get]
func handleGetCars(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}
		var filter model.Filter
		filter.Limit = uint(limit)
		filter.Offset = uint(offset)
		filter.RegNum, filter.ByRegNum = c.GetQuery("regNum")
		filter.Mark, filter.ByMark = c.GetQuery("mark")
		filter.Model, filter.ByModel = c.GetQuery("model")
		filter.OwnerName, filter.ByOwnerName = c.GetQuery("ownerName")
		filter.OwnerSurname, filter.ByOwnerSurname = c.GetQuery("ownerSurname")
		filter.OwnerPatronymic, filter.ByOwnerPatronymic = c.GetQuery("ownerPatronymic")
		if _, filter.ByYear = c.GetQuery("year"); filter.ByYear {
			filter.Year, err = strconv.Atoi(c.Query("year"))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
				return
			}
		}

		cars, err := a.GetCars(c, filter)

		switch {
		case err == nil:
			data := carsToCarsData(cars)
			c.JSON(http.StatusOK, carsResponse{
				Data: data,
				Err:  nil,
			})
		case errors.Is(err, model.ErrDatabaseError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrDatabaseError))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrServiceError))
		}
	}
}

// @Summary Добавление новых автомобилей
// @Description Получает на вход номера автомобилей, выполняет запрос во внешний API для получения недостающих данных и добавляет информацию о новых автомобилях
// @Accept json
// @Produce json
// @Param input body addCarsRequest true "Регистрационные номера"
// @Success 200 {object} carsResponse "Успешное добавление информации"
// @Failure 400 {object} carsResponse "Неверный формат входных данных"
// @Failure 409 {object} carsResponse "Попытка добавления существующего номера"
// @Failure 500 {object} carsResponse "Ошибка на стороне сервера"
// @Router /cars [post]
func handleAddCars(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req addCarsRequest
		if err := c.BindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}

		cars, err := a.AddCars(c, req.RegNums)

		switch {
		case err == nil:
			data := carsToCarsData(cars)
			c.JSON(http.StatusOK, carsResponse{
				Data: data,
				Err:  nil,
			})
		case errors.Is(err, model.ErrDatabaseError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrDatabaseError))
		case errors.Is(err, model.ErrDuplicateRegNum):
			c.AbortWithStatusJSON(http.StatusConflict, errorResponse(model.ErrDuplicateRegNum))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrServiceError))
		}
	}
}

// @Summary Изменение информации об автомобиле
// @Description Изменяет информацию о существующем автомобиле
// @Accept json
// @Produce json
// @Param id path int true "id автомобиля"
// @Param input body carData true "Новая информация об автомобиле"
// @Success 200 {object} carResponse "Успешное изменение информации"
// @Failure 400 {object} carResponse "Неверный формат входных данных"
// @Failure 404 {object} carResponse "Автомобиль с указанным id не найден"
// @Failure 409 {object} carResponse "Попытка добавления существующего номера"
// @Failure 500 {object} carResponse "Ошибка на стороне сервера"
// @Router /cars/:id [put]
func handleUpdateCar(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}
		var req carData
		if err = c.BindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}

		car, err := a.UpdateCar(c, id, model.Car{
			RegNum: req.RegNum,
			Mark:   req.Mark,
			Model:  req.Model,
			Year:   req.Year,
			Owner: model.Owner{
				Name:       req.Owner.Name,
				Surname:    req.Owner.Surname,
				Patronymic: req.Owner.Patronymic,
			},
		})

		switch {
		case err == nil:
			data := carToCarData(car)
			c.JSON(http.StatusOK, carResponse{
				Data: &data,
				Err:  nil,
			})
		case errors.Is(err, model.ErrCarNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse(model.ErrCarNotFound))
		case errors.Is(err, model.ErrDatabaseError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrDatabaseError))
		case errors.Is(err, model.ErrDuplicateRegNum):
			c.AbortWithStatusJSON(http.StatusConflict, errorResponse(model.ErrDuplicateRegNum))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrServiceError))
		}
	}
}

// @Summary Удаление информации об автомобиле
// @Description Удаляет информацию об автомобиле по его id
// @Produce json
// @Param id path int true "id автомобиля"
// @Success 200 {object} carResponse "Успешное удаление информации"
// @Failure 400 {object} carResponse "Неверный формат входных данных"
// @Failure 404 {object} carResponse "Автомобиль с указанным id не найден"
// @Failure 500 {object} carResponse "Ошибка на стороне сервера"
// @Router /cars/:id [delete]
func handleDeleteCar(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(model.ErrInvalidInput))
			return
		}

		err = a.DeleteCar(c, id)

		switch {
		case err == nil:
			c.AbortWithStatusJSON(http.StatusOK, errorResponse(nil))
		case errors.Is(err, model.ErrCarNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, errorResponse(model.ErrCarNotFound))
		case errors.Is(err, model.ErrDatabaseError):
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrDatabaseError))
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(model.ErrServiceError))
		}
	}
}

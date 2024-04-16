package api

import (
	"cars-service/internal/model"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type apiImpl struct {
	url string
	http.Client
}

func (a *apiImpl) GetInfo(ctx context.Context, regNum string) (model.Car, error) {
	params := url.Values{}
	params.Add("regNum", regNum)
	reqUrl := a.url + "?" + params.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	resp, err := a.Client.Do(req)
	if err != nil {
		return model.Car{}, errors.Join(model.ErrApiError, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return model.Car{}, model.ErrApiError
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Car{}, errors.Join(model.ErrApiError, err)
	}

	var data carInfo
	if err = json.Unmarshal(body, &data); err != nil {
		return model.Car{}, errors.Join(model.ErrApiError, err)
	}

	return model.Car{
		RegNum: data.RegNum,
		Mark:   data.Mark,
		Model:  data.Model,
		Year:   data.Year,
		Owner: model.Owner{
			Name:       data.Owner.Name,
			Surname:    data.Owner.Surname,
			Patronymic: data.Owner.Patronymic,
		},
	}, nil

}

// carInfo is a struct for parsing data from response body
type carInfo struct {
	RegNum string    `json:"regNum"`
	Mark   string    `json:"mark"`
	Model  string    `json:"model"`
	Year   int       `json:"year"`
	Owner  ownerInfo `json:"owner"`
}

// ownerInfo is a struct for parsing data from response body
type ownerInfo struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

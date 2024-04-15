package httpserver

type addCarsRequest struct {
	RegNums []string `json:"regNums"`
}

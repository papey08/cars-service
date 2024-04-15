package model

type Car struct {
	Id     uint64
	RegNum string
	Mark   string
	Model  string
	Year   int
	Owner
}

type Owner struct {
	Id         uint64
	Name       string
	Surname    string
	Patronymic string
}

type Filter struct {
	Limit  uint
	Offset uint

	RegNum          string
	Mark            string
	Model           string
	Year            int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string

	ByRegNum          bool
	ByMark            bool
	ByModel           bool
	ByYear            bool
	ByOwnerName       bool
	ByOwnerSurname    bool
	ByOwnerPatronymic bool
}

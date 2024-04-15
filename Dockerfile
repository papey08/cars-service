FROM golang:1.21.6

ENV GOPATH=/

WORKDIR /go/src/cars-service
COPY . .

RUN go mod download
RUN go build -o cars-service-app cmd/server/main.go

CMD ["./cars-service-app --docker"]

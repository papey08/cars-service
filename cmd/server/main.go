package main

import (
	"cars-service/internal/adapters/api"
	"cars-service/internal/app"
	"cars-service/internal/ports/httpserver"
	"cars-service/internal/repo"
	"cars-service/pkg/logger"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func connectToPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	coreRepoUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_DB_USERNAME"),
		os.Getenv("POSTGRES_DB_PASSWORD"),
		os.Getenv("POSTGRES_DB_HOST"),
		os.Getenv("POSTGRES_DB_PORT"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_DB_SSLMODE"))

	// 30 attempts to connect to postgres starting in docker container
	for i := 0; i < 30; i++ {
		conn, err := pgxpool.New(ctx, coreRepoUrl)
		if err != nil {
			time.Sleep(time.Second)
		} else {
			return conn, nil
		}
	}

	return nil, errors.New("unable to connect to postgres passwords repo")
}

const (
	dockerConfigFile = "config/docker/.env"
	localConfigFile  = "config/local/.env"
)

//	@title			cars-service API
//	@version		1.0
//	@description	Swagger-документация к API каталога автомобилей
//	@host			localhost:8080
//	@BasePath		/api/v1

func main() {
	logs := logger.New()

	isDocker := flag.Bool("docker", false, "flag if this project is running in docker container")
	flag.Parse()
	var configPath string
	if *isDocker {
		configPath = dockerConfigFile
	} else {
		configPath = localConfigFile
	}

	if err := godotenv.Load(configPath); err != nil {
		logs.Fatal(nil, "unable to load config files")
	}

	pool, err := connectToPostgres(context.Background())
	if err != nil {
		logs.Fatal(nil, err.Error())
	}
	defer pool.Close()

	apiCli := api.New(os.Getenv("API_ADDR"))

	a := app.New(repo.New(pool), apiCli, logs)

	srv := httpserver.New(os.Getenv("SERVER_ADDR"), a, logs)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logs.Fatal(nil, err.Error())
		}
	}()
	logs.Info(nil, "server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	<-quit
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = srv.Shutdown(shutdownCtx)
}

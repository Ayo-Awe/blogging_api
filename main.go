package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/ayo-awe/blogging_api/api"
	"github.com/ayo-awe/blogging_api/database"
	_ "github.com/ayo-awe/blogging_api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	PORT         int    `envconfig:"PORT" default:"8080"`
	DATABASE_URL string `envconfig:"DB_URL" required:"true"`
}

//	@title			Golang Blogging API
//	@version		1.0
//	@description	This is a minimalist blogging api.

// @BasePath	/api
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	r := chi.NewRouter()
	logger := slog.Default()
	cfg, err := LoadConfig()

	if err != nil {
		return err
	}

	db, err := database.NewDatabase(cfg.DATABASE_URL)
	if err != nil {
		return err
	}

	repo := database.NewArticleRepository(db)
	app := api.NewApplication(logger, repo)

	r.Use(middleware.Logger)
	r.Mount("/api", app.BuildRoutes())
	r.Get("/swagger/*", httpSwagger.Handler())

	fmt.Println("starting server on port 8080")
	addr := fmt.Sprintf(":%d", cfg.PORT)
	if err = http.ListenAndServe(addr, r); err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

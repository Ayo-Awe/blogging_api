package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ayo-awe/blogging_api/api"
	"github.com/ayo-awe/blogging_api/database"
	_ "github.com/ayo-awe/blogging_api/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Golang Blogging API
//	@version		1.0
//	@description	This is a minimalist blogging api.

//	@BasePath	/api
func main() {
	r := chi.NewRouter()
	logger := slog.Default()

	// todo: fetch db_url and port from env
	db, err := database.NewDatabase("postgresql://aweayo:aweayo@localhost:5432/blogging_api?sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	repo := database.NewArticleRepository(db)
	app := api.NewApplication(logger, repo)

	r.Use(middleware.Logger)
	r.Mount("/api", app.BuildRoutes())
	r.Get("/swagger/*", httpSwagger.Handler())

	fmt.Println("starting server on port 8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		logger.Error(err.Error())
	}
}

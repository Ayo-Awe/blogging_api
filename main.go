package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ayo-awe/blogging_api/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	logger := slog.Default()

	db, err := database.NewDatabase("postgresql://aweayo:aweayo@localhost:5432/blogging_api?sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	repo := database.NewArticleRepository(db)
	app := Application{repo: repo, logger: logger}

	r.Use(middleware.Logger)
	r.Mount("/api", app.BuildRoutes())

	fmt.Println("starting server on port 8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		logger.Error(err.Error())
	}
}

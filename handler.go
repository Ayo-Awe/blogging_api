package main

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/ayo-awe/blogging_api/database"
	"github.com/ayo-awe/blogging_api/utils"
	"github.com/go-chi/chi/v5"
)

type Application struct {
	logger *slog.Logger
	repo   database.ArticleRepository
}

type envelope map[string]interface{}

func (a *Application) BuildRoutes() chi.Router {
	router := chi.NewRouter()
	router.Route("/articles", func(r chi.Router) {
		r.Post("/", a.CreateArticle)
		r.Get("/", a.GetArticles)
		r.Get("/{id}", a.GetArticleByID)
	})

	return router
}

func (a *Application) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var payload database.Article

	err := utils.DecodeJSON(r, &payload)
	if err != nil {
		msg := "Please provide a valid JSON body"
		utils.HTTPError(w, http.StatusBadRequest, msg)
		return
	}

	err = payload.Validate()
	if err != nil {
		msg := err.Error()
		utils.HTTPError(w, http.StatusBadRequest, msg)
		return
	}

	article, err := a.repo.CreateArticle(r.Context(), &payload)
	if err != nil {
		msg := "An unexpected error occured"
		utils.HTTPError(w, http.StatusInternalServerError, msg)
		a.logger.Error(err.Error())
		return
	}

	utils.Created(w, envelope{"article": article}, nil)
}

func (a *Application) GetArticles(w http.ResponseWriter, r *http.Request) {
	// get rawTags  from query params
	var tags database.Tags

	rawTags := r.URL.Query().Get("tags")
	if rawTags != "" {
		mapFn := func(ele string) string { return strings.ToLower(strings.TrimSpace(ele)) }
		tags = utils.Map(strings.Split(rawTags, ","), mapFn)
	}

	// get articles by tags
	articles, err := a.repo.GetArticles(r.Context(), database.ArticleFilter{Tags: tags})
	if err != nil {
		utils.HTTPError(w, http.StatusInternalServerError, "An unexpected error occured")
		a.logger.Error(err.Error())
		return
	}

	utils.Ok(w, envelope{"articles": articles}, nil)
}

func (a *Application) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.HTTPError(w, http.StatusNotFound, "Article not found")
		return
	}

	article, err := a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.HTTPError(w, http.StatusNotFound, "Article not found")
			return
		}

		utils.HTTPError(w, http.StatusInternalServerError, "An unexpected error occured")
		a.logger.Error("GetArticleByID: " + err.Error())
		return
	}

	utils.Ok(w, envelope{"article": article}, nil)
}

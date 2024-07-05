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
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Application struct {
	logger *slog.Logger
	repo   database.ArticleRepository
}

type UpdateArticleRequest struct {
	Title   string
	Content string
	Tags    database.Tags
}

func (u *UpdateArticleRequest) Validate() error {
	u.clean()
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Length(5, 255)),
		validation.Field(&u.Content, validation.Length(5, 0)),
		validation.Field(&u.Tags, validation.Each(validation.Length(2, 0), is.LowerCase)),
	)
}

func (u *UpdateArticleRequest) clean() {
	u.Title = strings.TrimSpace(u.Title)
	u.Content = strings.TrimSpace(u.Content)

	for i, tag := range u.Tags {
		trimmed := strings.TrimSpace(tag)
		u.Tags[i] = strings.ToLower(trimmed)
	}
}

type envelope map[string]interface{}

func (a *Application) BuildRoutes() chi.Router {
	router := chi.NewRouter()
	router.Route("/articles", func(r chi.Router) {
		r.Post("/", a.CreateArticle)
		r.Get("/", a.GetArticles)
		r.Get("/{id}", a.GetArticleByID)
		r.Patch("/{id}", a.UpdateArticle)
		r.Delete("/{id}", a.DeleteArticle)
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
	// todo: paginate articles
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

func (a *Application) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.HTTPError(w, http.StatusNotFound, "Article Not Found")
		return
	}

	// find article by id
	article, err := a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.HTTPError(w, http.StatusNotFound, "Article Not Found")
			return
		}
		a.logger.Error("Update Article: " + err.Error())
		utils.HTTPError(w, http.StatusInternalServerError, "An unexpected error occured")
		return
	}

	// parse request body
	var payload UpdateArticleRequest
	err = utils.DecodeJSON(r, &payload)
	if err != nil {
		a.logger.Error("Update Article: " + err.Error())
		utils.HTTPError(w, http.StatusInternalServerError, "Failed to decode JSON request body")
		return
	}

	// validate request body
	if err = payload.Validate(); err != nil {
		utils.HTTPError(w, http.StatusBadRequest, err.Error())
		return
	}

	// update article with changes from request body
	if payload.Title != "" {
		article.Title = payload.Title
	}

	if payload.Content != "" {
		article.Content = payload.Content
	}

	if len(payload.Tags) > 0 {
		article.Tags = payload.Tags
	}

	// save updates in the databse
	updatedArticle, err := a.repo.UpdateArticle(r.Context(), article)
	if err != nil {
		a.logger.Error("Update Article: " + err.Error())
		utils.HTTPError(w, http.StatusInternalServerError, "Failed to decode JSON request body")
		return
	}

	utils.Ok(w, envelope{"article": updatedArticle}, nil)
}

func (a *Application) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// get article id
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.HTTPError(w, http.StatusNotFound, "Article Not Found")
		return
	}

	_, err = a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.HTTPError(w, http.StatusNotFound, "Article Not Found")
			return
		}
		a.logger.Error("Delete Article" + err.Error())
		utils.HTTPError(w, http.StatusInternalServerError, "An unexpected error occured")
		return
	}

	err = a.repo.DeleteArticle(r.Context(), id)
	if err != nil {

		a.logger.Error("Delete Article" + err.Error())
		utils.HTTPError(w, http.StatusInternalServerError, "An unexpected error occured")
		return
	}

	utils.NoContent(w)
}

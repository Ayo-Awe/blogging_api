package api

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

const (
	DEFAULT_PAGE     = 1
	DEFAULT_PER_PAGE = 20
	MAX_PER_PAGE     = 100
)

type Application struct {
	logger *slog.Logger
	repo   database.ArticleRepository
}

func NewApplication(logger *slog.Logger, repo database.ArticleRepository) *Application {
	return &Application{logger, repo}
}

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

// CreateArticle godoc
//	@Summary	Create article
//	@Tags		articles
//	@Accept		json
//	@Produce	json
//	@Param		data	body		CreateArticleRequest	true	"Request Body"
//	@Success	200		{object}	SuccessReponse{data=CreateArticleResponse}
//	@Failure	400		{object}	ErrorResponse
//	@Router		/articles [post]
func (a *Application) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var payload CreateArticleRequest

	err := utils.DecodeJSON(r, &payload)
	if err != nil {
		msg := "Please provide a valid JSON body"
		utils.RenderResponse(w, http.StatusBadRequest, NewErrResponse(msg))
		return
	}

	err = payload.Validate()
	if err != nil {
		msg := err.Error()
		utils.RenderResponse(w, http.StatusBadRequest, NewErrResponse(msg))
		return
	}

	article, err := a.repo.CreateArticle(r.Context(), payload.toArticle())
	if err != nil {
		msg := "An unexpected error occured"
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse(msg))
		a.logger.Error(err.Error())
		return
	}

	data := CreateArticleResponse{Article: *article}
	utils.RenderResponse(w, http.StatusCreated, NewSuccessResponse(data, nil))
}

// GetArticles godoc
//	@Summary	List article
//	@Tags		articles
//	@Accept		json
//	@Produce	json
//	@Param		tags	query		[]string	false	"Filter by tags"
//	@Param		page	query		int			false	"Page"
//	@Param		perPage	query		int			false	"Articles per page"
//	@Success	200		{object}	SuccessReponse{data=GetArticlesResponse,metadata=database.PaginationData}
//	@Router		/articles [get]
func (a *Application) GetArticles(w http.ResponseWriter, r *http.Request) {
	// get rawTags  from query params
	var tags database.Tags

	rawTags := r.URL.Query().Get("tags")
	if rawTags != "" {
		mapFn := func(ele string) string { return strings.ToLower(strings.TrimSpace(ele)) }
		tags = utils.Map(strings.Split(rawTags, ","), mapFn)
	}

	rawPage := r.URL.Query().Get("page")
	page, err := strconv.Atoi(rawPage)
	if err != nil || page <= 0 {
		page = DEFAULT_PAGE
	}

	rawPerPage := r.URL.Query().Get("per_page")
	perPage, err := strconv.Atoi(rawPerPage)
	if err != nil {
		perPage = DEFAULT_PER_PAGE
	}

	if perPage > MAX_PER_PAGE || perPage <= 0 {
		perPage = MAX_PER_PAGE
	}

	pageable := database.Paging{
		Page:    page,
		PerPage: perPage,
	}

	// get articles by tags
	articles, paginationData, err := a.repo.GetArticles(r.Context(), database.ArticleFilter{Tags: tags}, pageable)
	if err != nil {
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("An unexpected error occured"))
		a.logger.Error(err.Error())
		return
	}

	data := GetArticlesResponse{Articles: articles}
	utils.RenderResponse(w, http.StatusOK, NewSuccessResponse(data, paginationData))
}

// GetArticleByID godoc
//	@Summary	Get article by ID
//	@Tags		articles
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"Article ID"
//	@Success	200	{object}	SuccessReponse{data=GetArticleByIDResponse}
//	@Failure	404	{object}	ErrorResponse
//	@Router		/articles/{id} [get]
func (a *Application) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article not found"))
		return
	}

	article, err := a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article not found"))
			return
		}

		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("An unexpected error occured"))
		a.logger.Error("GetArticleByID: " + err.Error())
		return
	}

	data := GetArticleByIDResponse{Article: *article}
	utils.RenderResponse(w, http.StatusOK, NewSuccessResponse(data, nil))
}

// UpdateArticle godoc
//	@Summary	Update article
//	@Tags		articles
//	@Accept		json
//	@Produce	json
//	@Param		id		path		int						true	"Article ID"
//	@Param		data	body		UpdateArticleRequest	true	"Request Body"
//	@Success	200		{object}	SuccessReponse{data=UpdateArticleResponse}
//	@Failure	400		{object}	ErrorResponse
//	@Failure	404		{object}	ErrorResponse
//	@Router		/articles/{id} [patch]
func (a *Application) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article Not Found"))
		return
	}

	// find article by id
	article, err := a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article Not Found"))
			return
		}
		a.logger.Error("Update Article: " + err.Error())
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("An unexpected error occured"))
		return
	}

	// parse request body
	var payload UpdateArticleRequest
	err = utils.DecodeJSON(r, &payload)
	if err != nil {
		a.logger.Error("Update Article: " + err.Error())
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("Failed to decode JSON request body"))
		return
	}

	// validate request body
	if err = payload.Validate(); err != nil {
		utils.RenderResponse(w, http.StatusBadRequest, NewErrResponse(err.Error()))
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
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("Failed to decode JSON request body"))
		return
	}

	data := UpdateArticleResponse{Article: *updatedArticle}
	utils.RenderResponse(w, http.StatusOK, NewSuccessResponse(data, nil))
}

// DeleteArticle godoc
//	@Summary	Delete article
//	@Tags		articles
//	@Accept		json
//	@Produce	json
//	@Param		id	path	int	true	"Article ID"
//	@Success	204
//	@Failure	404	{object}	ErrorResponse
//	@Router		/articles/{id} [delete]
func (a *Application) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// get article id
	rawID := chi.URLParam(r, "id")

	id, err := strconv.Atoi(rawID)
	if err != nil {
		utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article Not Found"))
		return
	}

	_, err = a.repo.GetArticleByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrArticleNotFound) {
			utils.RenderResponse(w, http.StatusNotFound, NewErrResponse("Article Not Found"))
			return
		}
		a.logger.Error("Delete Article" + err.Error())
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("An unexpected error occured"))
		return
	}

	err = a.repo.DeleteArticle(r.Context(), id)
	if err != nil {

		a.logger.Error("Delete Article" + err.Error())
		utils.RenderResponse(w, http.StatusInternalServerError, NewErrResponse("An unexpected error occured"))
		return
	}

	utils.RenderResponse(w, http.StatusNoContent, nil)
}

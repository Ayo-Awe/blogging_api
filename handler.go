package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ayo-awe/blogging_api/database"
)

type Application struct {
	logger *slog.Logger
	repo   database.ArticleRepository
}

type response struct {
	Status   string      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Message  string      `json:"message,omitempty"`
}

func (a *Application) ok(w http.ResponseWriter, data interface{}, metadata interface{}) {
	res := response{
		Data:     data,
		Status:   "success",
		Metadata: metadata,
	}
	a.JSONResponse(w, http.StatusOK, res)
}

func (a *Application) created(w http.ResponseWriter, data interface{}, metadata interface{}) {
	res := response{
		Data:     data,
		Status:   "success",
		Metadata: metadata,
	}
	a.JSONResponse(w, http.StatusOK, res)

}

func (a *Application) JSONResponse(w http.ResponseWriter, statusCode int, response response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (a *Application) httpError(w http.ResponseWriter, statusCode int, msg string) {
	res := response{
		Status:  "error",
		Message: msg,
	}
	a.JSONResponse(w, statusCode, res)
}

func (a *Application) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var payload database.Article

	err := decodeJSON(r, &payload)
	if err != nil {
		msg := "Please provide a valid JSON body"
		a.httpError(w, http.StatusBadRequest, msg)
		return
	}

	err = payload.Validate()
	if err != nil {
		msg := err.Error()
		a.httpError(w, http.StatusBadRequest, msg)
		return
	}

	article, err := a.repo.CreateArticle(r.Context(), &payload)
	if err != nil {
		msg := "An unexpected error occured"
		a.httpError(w, http.StatusInternalServerError, msg)
		a.logger.Error(err.Error())
		return
	}

	data := map[string]interface{}{"article": article}
	a.created(w, data, nil)
}

func decodeJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

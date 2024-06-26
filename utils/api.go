package utils

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status   string      `json:"status"`
	Data     interface{} `json:"data,omitempty"`
	Metadata interface{} `json:"metadata,omitempty"`
	Message  string      `json:"message,omitempty"`
}

func Ok(w http.ResponseWriter, data interface{}, metadata interface{}) {
	res := response{
		Data:     data,
		Status:   "success",
		Metadata: metadata,
	}
	JSONResponse(w, http.StatusOK, res)
}

func Created(w http.ResponseWriter, data interface{}, metadata interface{}) {
	res := response{
		Data:     data,
		Status:   "success",
		Metadata: metadata,
	}
	JSONResponse(w, http.StatusOK, res)

}

func JSONResponse(w http.ResponseWriter, statusCode int, response response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func HTTPError(w http.ResponseWriter, statusCode int, msg string) {
	res := response{
		Status:  "error",
		Message: msg,
	}
	JSONResponse(w, statusCode, res)
}

func DecodeJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dest)
}

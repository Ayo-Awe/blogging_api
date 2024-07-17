package utils

import (
	"encoding/json"
	"net/http"
)

func RenderResponse(w http.ResponseWriter, statusCode int, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if res != nil {
		json.NewEncoder(w).Encode(res)
	}
}

func DecodeJSON(r *http.Request, dest interface{}) error {
	decoder := json.NewDecoder(r.Body)

	return decoder.Decode(dest)
}

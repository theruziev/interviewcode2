package httpx

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

func ParseJSON[T any](r *http.Request) (*T, error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func JSONResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		jsonMarshalErr := fmt.Errorf("failed to marshal json: %w", err)
		JSONError(w, http.StatusInternalServerError, jsonMarshalErr.Error())
		return
	}
}

func JSONOKResponse(w http.ResponseWriter) {
	JSONResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func JSONError(w http.ResponseWriter, code int, text string) {
	data := map[string]string{
		"error": text,
	}
	JSONResponse(w, code, data)
}

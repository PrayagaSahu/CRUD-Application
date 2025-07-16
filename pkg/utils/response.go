package utils

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := AppError{
		Code:    statusCode,
		Message: http.StatusText(statusCode),
	}
	json.NewEncoder(w).Encode(err)
}

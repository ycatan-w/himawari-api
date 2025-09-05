package utils

import (
	"encoding/json"
	"net/http"
)

// API JSON responses struct
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []APIError  `json:"errors,omitempty"`
	Message string      `json:"message,omitempty"`
}

type APIError struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// API response helpers
func RespondJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func CreateResponseError(message string, errors []APIError) APIResponse {
	return APIResponse{
		Success: false,
		Errors:  errors,
		Message: message,
	}
}

func BuildApiError(field string, code string, details string) APIError {
	return APIError{Field: field, Code: code, Details: details}
}

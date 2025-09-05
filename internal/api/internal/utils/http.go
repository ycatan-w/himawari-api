package utils

import (
	"fmt"
	"net/http"
)

func ValidateGetRequest(w http.ResponseWriter, method string) bool {
	return validateRequestMethod(w, http.MethodGet, method)
}

func ValidatePostRequest(w http.ResponseWriter, method string) bool {
	return validateRequestMethod(w, http.MethodPost, method)
}

func ValidatePutRequest(w http.ResponseWriter, method string) bool {
	return validateRequestMethod(w, http.MethodPut, method)
}

func ValidateDeleteRequest(w http.ResponseWriter, method string) bool {
	return validateRequestMethod(w, http.MethodDelete, method)
}

func validateRequestMethod(w http.ResponseWriter, expectedMethod string, currentMethod string) bool {
	if currentMethod != expectedMethod {
		RespondJSON(w, http.StatusMethodNotAllowed, CreateResponseError(
			"method_not_allowed",
			[]APIError{BuildApiError("", "METHOD_NOT_ALLOWED", fmt.Sprintf("Expected Method \"%s\" but got \"%s\"", expectedMethod, currentMethod))},
		))

		return false
	}

	return true
}

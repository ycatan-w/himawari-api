package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ycatan-w/himawari-api/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GenerateAndSaveToken(w http.ResponseWriter, userId int) (string, bool) {
	token, err := generateToken()
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, CreateResponseError("token_failed", []APIError{BuildApiError("", "TOKEN_FAILED", err.Error())}))
		return "", false
	}

	db.AddUserToken(userId, token)
	return token, true
}

func EncryptPassword(w http.ResponseWriter, password string) (string, bool) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		RespondJSON(w, http.StatusBadRequest, CreateResponseError("hash_failed", []APIError{BuildApiError("", "HASH_FAILED", err.Error())}))

		return "", false
	}

	return string(hash), true
}

func ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func GetFormData[T any](w http.ResponseWriter, body io.ReadCloser) *T {
	var req T
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		RespondJSON(w, http.StatusBadRequest, CreateResponseError(
			"invalid_request",
			[]APIError{BuildApiError("", "INVALID_REQUEST", "Invalid JSON")},
		))
		return nil
	}

	return &req
}

func GetIdFromPath(r *http.Request, param string) (int, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		return 0, errors.New(param + " is required")
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid " + param)
	}

	return id, nil
}

package api

import (
	"encoding/json"
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/login"
	"github.com/ycatan-w/himawari-api/internal/api/internal/register"
	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
)

// POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePostRequest(w, r.Method) {
		return
	}

	reqData := register.GetRequestUser(w, r.Body)
	if reqData == nil || !register.ValidateUser(w, reqData) {
		return
	}

	token, ok := register.Register(w, reqData)
	if !ok {
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.APIResponse{
		Success: true,
		Data:    map[string]string{"username": reqData.Username, "token": token},
		Message: "user_created",
	})
}

// POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePostRequest(w, r.Method) {
		return
	}

	reqData := login.GetRequestUser(w, r.Body)
	if reqData == nil || !login.ValidateUser(w, reqData) {
		return
	}

	token, ok := login.Login(w, reqData)
	if !ok {
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.APIResponse{
		Success: true,
		Data:    map[string]string{"username": reqData.Username, "token": token},
		Message: "user_found",
	})
}

// POST /logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePostRequest(w, r.Method) {
		return
	}
	type TokenFormData struct {
		Token string `json:"token"`
	}
	var req TokenFormData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_request",
			[]utils.APIError{utils.BuildApiError("", "INVALID_REQUEST", "Invalid JSON")},
		))
		return
	}
	db.RemoveUserToken(req.Token)

	w.WriteHeader(http.StatusNoContent)
}

package api

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/log"
	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
	"github.com/ycatan-w/himawari-api/internal/db"
)

// -------------------- Server Handlers --------------------
func LogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /logs
		handleGetLogs(w, r)
	case http.MethodPost:
		// POST /logs
		handlePostLog(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func LogsByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		// PUT /logs/{id}
		handlePutLog(w, r)
	case http.MethodDelete:
		// DELETE /logs/{id}
		handleDeleteLog(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// -------------------- /logs --------------------
func handleGetLogs(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateGetRequest(w, r.Method) {
		return
	}
	var endpoint = log.GetLogEndPoint{}
	userId, date, ok := endpoint.ValidateGet(w, r)
	if !ok {
		return
	}
	logs, err := db.FindLogsByUserAndDate(userId, date)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_query_failed", []utils.APIError{utils.BuildApiError("", "DB_QUERY_FAILED", err.Error())}))
		return
	}

	utils.RespondJSON(w, http.StatusOK,
		utils.APIResponse{
			Success: true,
			Data:    logs,
			Message: "log_found",
		})
}
func handlePostLog(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePostRequest(w, r.Method) {
		return
	}
	userId, ok := middleware.GetUserID(r)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_user_id",
			[]utils.APIError{utils.BuildApiError("", "INVALID_USER_ID", "User id not found")},
		))
		return
	}
	var endpoint = log.PostLogEndPoint{}
	reqData := endpoint.GetLogFormData(w, r.Body)
	if reqData == nil || !endpoint.ValidateLogFormData(w, reqData) {
		return
	}

	log := endpoint.CreateLog(w, reqData, userId)
	if log == nil {
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.APIResponse{
		Success: true,
		Data:    log,
		Message: "log_created",
	})
}

// -------------------- /logs/{id} --------------------
func handlePutLog(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePutRequest(w, r.Method) {
		return
	}
	var endpoint = log.PutLogEndPoint{}
	reqData := endpoint.GetLogFormData(w, r)
	if reqData == nil || !endpoint.ValidateLogFormData(w, reqData) {
		return
	}

	log := endpoint.EditLog(w, reqData)
	if log == nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"log_updated_not_found",
			[]utils.APIError{utils.BuildApiError("", "LOG_UPDATED_NOT_FOUND", "Unable to find the updated log")},
		))
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.APIResponse{
		Success: true,
		Data:    log,
		Message: "log_updated",
	})
}
func handleDeleteLog(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateDeleteRequest(w, r.Method) {
		return
	}
	userId, ok := middleware.GetUserID(r)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_user_id",
			[]utils.APIError{utils.BuildApiError("", "INVALID_USER_ID", "User id not found")},
		))
		return
	}
	var endpoint = log.DeleteLogEndPoint{}
	logId := endpoint.GetLogId(w, r)
	if logId == nil {
		return
	}

	db.RemoveLog(*logId, userId)

	w.WriteHeader(http.StatusNoContent)
}

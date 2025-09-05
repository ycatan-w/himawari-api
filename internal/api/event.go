package api

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/event"
	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
	"github.com/ycatan-w/himawari-api/internal/db"
)

// -------------------- Server Handlers --------------------
func EventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET /events
		handleGetEvents(w, r)
	case http.MethodPost:
		// POST /events
		handlePostEvent(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func EventsByIdHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		// PUT /events/{id}
		handlePutEvent(w, r)
	case http.MethodDelete:
		// DELETE /events/{id}
		handleDeleteEvent(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// -------------------- /events --------------------
func handleGetEvents(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateGetRequest(w, r.Method) {
		return
	}
	var endpoint = event.GetEventEndPoint{}
	userId, date, ok := endpoint.ValidateGet(w, r)
	if !ok {
		return
	}
	events, err := db.FindEventsByUserAndDate(userId, date)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_query_failed", []utils.APIError{utils.BuildApiError("", "DB_QUERY_FAILED", err.Error())}))
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.APIResponse{
		Success: true,
		Data:    events,
		Message: "event_found",
	})
}

func handlePostEvent(w http.ResponseWriter, r *http.Request) {
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
	var endpoint = event.PostEventEndPoint{}
	reqData := endpoint.GetEventFormData(w, r.Body)
	if reqData == nil || !endpoint.ValidateEventFormData(w, reqData) {
		return
	}

	event := endpoint.CreateEvent(w, reqData, userId)
	if event == nil {
		return
	}

	utils.RespondJSON(w, http.StatusCreated, utils.APIResponse{
		Success: true,
		Data:    event,
		Message: "event_created",
	})
}

// -------------------- /events/{id} --------------------
func handlePutEvent(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidatePutRequest(w, r.Method) {
		return
	}
	var endpoint = event.PutEventEndPoint{}
	reqData := endpoint.GetEventFormData(w, r)
	if reqData == nil || !endpoint.ValidateEventFormData(w, reqData) {
		return
	}

	event := endpoint.EditEvent(w, reqData)
	if event == nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"event_updated_not_found",
			[]utils.APIError{utils.BuildApiError("", "EVENT_UPDATED_NOT_FOUND", "Unable to find the updated event")},
		))
		return
	}

	utils.RespondJSON(w, http.StatusOK, utils.APIResponse{
		Success: true,
		Data:    event,
		Message: "event_updated",
	})
}

func handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
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
	var endpoint = event.DeleteEventEndPoint{}
	eventId := endpoint.GetEventId(w, r)
	if eventId == nil {
		return
	}

	db.RemoveEvent(*eventId, userId)

	w.WriteHeader(http.StatusNoContent)
}

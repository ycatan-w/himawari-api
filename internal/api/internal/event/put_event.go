package event

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
	"github.com/ycatan-w/himawari-api/internal/db"
)

type PutEventEndPoint struct{}

func (PutEventEndPoint) GetEventFormData(w http.ResponseWriter, r *http.Request) *EventFormData {
	requestEvent := utils.GetFormData[EventFormData](w, r.Body)
	if requestEvent == nil {
		return nil
	}
	userId, ok := middleware.GetUserID(r)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_user_id",
			[]utils.APIError{utils.BuildApiError("", "INVALID_USER_ID", "User id not found")},
		))
		return nil
	}
	requestEvent.UserId = userId

	idInt := GetEventId(w, r)
	if idInt == nil {
		return nil
	}
	requestEvent.Id = *idInt

	return requestEvent
}

func (PutEventEndPoint) ValidateEventFormData(w http.ResponseWriter, requestEvent *EventFormData) bool {
	errs := CheckEventFormData(requestEvent, nil)

	if requestEvent.Id <= 0 {
		errs = append(errs, utils.BuildApiError("event_id", "REQUIRED_EVENT_ID", "EventId cannot be empty."))
	} else {
		event := db.FindOneEventByIdAndUserId(requestEvent.Id, requestEvent.UserId)
		if event == nil {
			errs = append(errs, utils.BuildApiError("event_id", "INVALID_EVENT_ID", "The event id provided wasn't found."))
		}
	}

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))

		return false
	}
	return true
}

func (PutEventEndPoint) EditEvent(w http.ResponseWriter, requestEvent *EventFormData) *db.Event {
	_, err := db.UpdateEvent(db.Event{
		Title:       requestEvent.Title,
		Description: requestEvent.Description,
		Date:        requestEvent.Date,
		Start:       requestEvent.Start,
		End:         requestEvent.End,
	}, requestEvent.Id, requestEvent.UserId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_connection_failed", []utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())}))
		return nil
	}
	return db.FindOneEventById(requestEvent.Id)
}

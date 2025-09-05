package event

import (
	"io"
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
)

type PostEventEndPoint struct{}

func (PostEventEndPoint) GetEventFormData(w http.ResponseWriter, body io.ReadCloser) *EventFormData {
	return utils.GetFormData[EventFormData](w, body)
}

func (PostEventEndPoint) ValidateEventFormData(w http.ResponseWriter, requestEvent *EventFormData) bool {
	errs := CheckEventFormData(requestEvent, nil)

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))

		return false
	}

	return true
}

func (PostEventEndPoint) CreateEvent(w http.ResponseWriter, requestEvent *EventFormData, userId int) *db.Event {
	createdId, err := db.AddEvent(db.Event{
		Title:       requestEvent.Title,
		Description: requestEvent.Description,
		Date:        requestEvent.Date,
		Start:       requestEvent.Start,
		End:         requestEvent.End,
	}, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_connection_failed", []utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())}))
		return nil
	}
	event := db.FindOneEventById(createdId)
	if event == nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"event_created_not_found",
			[]utils.APIError{utils.BuildApiError("", "EVENT_CREATED_NOT_FOUND", "Unable to find the created event")},
		))
		return nil
	}
	return event
}

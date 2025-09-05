package event

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
)

type EventFormData struct {
	Id          int
	Title       string
	Description string
	Date        string
	Start       int
	End         int
	UserId      int
}

func CheckEventFormData(requestEvent *EventFormData, errs []utils.APIError) []utils.APIError {
	if len(requestEvent.Title) <= 0 {
		errs = append(errs, utils.BuildApiError("title", "REQUIRED_TITLE", "Title cannot be empty."))
	}
	if len(requestEvent.Title) > 100 {
		errs = append(errs, utils.BuildApiError("title", "INVALID_TITLE_LENGTH", "Title cannot contains more than 100 chars."))
	}
	if len(requestEvent.Description) > 1000 {
		errs = append(errs, utils.BuildApiError("description", "INVALID_DESCRIPTION_LENGTH", "Description cannot contains more than 1000 chars."))
	}
	if !utils.ValidateDate(requestEvent.Date) {
		errs = append(errs, utils.BuildApiError("date", "INVALID_DATE_FORMAT", "Date must be in YYYY-MM-DD format."))
	}
	if requestEvent.Start < 0 || requestEvent.Start > 1439 {
		errs = append(errs, utils.BuildApiError("start", "INVALID_START", "Start must be between 0 and 1439 (00:00 and 23:59)"))
	}
	if requestEvent.End < 0 || requestEvent.End > 1439 {
		errs = append(errs, utils.BuildApiError("end", "INVALID_END", "End must be between 0 and 1439 (00:00 and 23:59)"))
	}
	if requestEvent.Start >= requestEvent.End {
		errs = append(errs, utils.BuildApiError("start", "INVALID_TIME_RANGE", "Start cannot be after End"))
	}

	return errs
}

func GetEventId(w http.ResponseWriter, r *http.Request) *int {
	id, err := utils.GetIdFromPath(r, "event_id")
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_event_id",
			[]utils.APIError{utils.BuildApiError("event_id", "INVALID_EVENT_ID", err.Error())},
		))
		return nil
	}
	return &id
}

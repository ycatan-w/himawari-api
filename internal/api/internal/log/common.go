package log

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
)

type LogFormData struct {
	Id     int
	Text   string
	Date   string
	UserId int
}

func CheckLogFormData(requestEvent *LogFormData, errs []utils.APIError) []utils.APIError {
	if len(requestEvent.Text) <= 0 {
		errs = append(errs, utils.BuildApiError("text", "REQUIRED_TEXT", "Text cannot be empty."))
	}
	if len(requestEvent.Text) > 5000 {
		errs = append(errs, utils.BuildApiError("text", "INVALID_TEXT_LENGTH", "Text cannot contains more than 5000 chars."))
	}
	if !utils.ValidateDate(requestEvent.Date) {
		errs = append(errs, utils.BuildApiError("date", "INVALID_DATE_FORMAT", "Date must be in YYYY-MM-DD format."))
	}

	return errs
}

func GetLogId(w http.ResponseWriter, r *http.Request) *int {
	id, err := utils.GetIdFromPath(r, "log_id")
	if err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_log_id",
			[]utils.APIError{utils.BuildApiError("log_id", "INVALID_LOG_ID", err.Error())},
		))
		return nil
	}
	return &id
}

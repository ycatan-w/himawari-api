package log

import (
	"io"
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/db"
)

type PostLogEndPoint struct{}

func (PostLogEndPoint) GetLogFormData(w http.ResponseWriter, body io.ReadCloser) *LogFormData {
	return utils.GetFormData[LogFormData](w, body)
}

func (PostLogEndPoint) ValidateLogFormData(w http.ResponseWriter, requestLog *LogFormData) bool {
	errs := CheckLogFormData(requestLog, nil)

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))

		return false
	}

	return true
}

func (PostLogEndPoint) CreateLog(w http.ResponseWriter, requestLog *LogFormData, userId int) *db.Log {
	createdId, err := db.AddLog(db.Log{
		Text: requestLog.Text,
		Date: requestLog.Date,
	}, userId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_connection_failed", []utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())}))
		return nil
	}
	log := db.FindOneLogById(createdId)
	if log == nil {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"log_created_not_found",
			[]utils.APIError{utils.BuildApiError("", "LOG_CREATED_NOT_FOUND", "Unable to find the created Log")},
		))
		return nil
	}
	return log
}

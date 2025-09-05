package log

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
	"github.com/ycatan-w/himawari-api/internal/db"
)

type PutLogEndPoint struct{}

func (PutLogEndPoint) GetLogFormData(w http.ResponseWriter, r *http.Request) *LogFormData {
	requestLog := utils.GetFormData[LogFormData](w, r.Body)
	if requestLog == nil {
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
	requestLog.UserId = userId

	idInt := GetLogId(w, r)
	if idInt == nil {
		return nil
	}
	requestLog.Id = *idInt

	return requestLog
}

func (PutLogEndPoint) ValidateLogFormData(w http.ResponseWriter, requestLog *LogFormData) bool {
	errs := CheckLogFormData(requestLog, nil)

	if requestLog.Id <= 0 {
		errs = append(errs, utils.BuildApiError("log_id", "REQUIRED_LOG_ID", "LogId cannot be empty."))
	} else {
		event := db.FindOneLogByIdAndUserId(requestLog.Id, requestLog.UserId)
		if event == nil {
			errs = append(errs, utils.BuildApiError("log_id", "INVALID_LOG_ID", "The log id provided wasn't found."))
		}
	}

	if len(errs) > 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError("validation_failed", errs))

		return false
	}
	return true
}

func (PutLogEndPoint) EditLog(w http.ResponseWriter, requestLog *LogFormData) *db.Log {
	_, err := db.UpdateLog(db.Log{
		Text: requestLog.Text,
		Date: requestLog.Date,
	}, requestLog.Id, requestLog.UserId)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, utils.CreateResponseError("db_connection_failed", []utils.APIError{utils.BuildApiError("", "DB_CONNECTION_FAILED", err.Error())}))
		return nil
	}
	return db.FindOneLogById(requestLog.Id)
}

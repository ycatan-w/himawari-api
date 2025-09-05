package event

import (
	"net/http"

	"github.com/ycatan-w/himawari-api/internal/api/internal/utils"
	"github.com/ycatan-w/himawari-api/internal/api/middleware"
)

type GetEventEndPoint struct{}

func (GetEventEndPoint) ValidateGet(w http.ResponseWriter, r *http.Request) (int, string, bool) {
	userId, ok := middleware.GetUserID(r)
	if !ok {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_user_id",
			[]utils.APIError{utils.BuildApiError("", "INVALID_USER_ID", "User id not found")},
		))
		return -1, "", false
	}
	date := r.URL.Query().Get("date") // "YYYY-MM-DD"
	if len(date) <= 0 {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"required_date",
			[]utils.APIError{utils.BuildApiError("date", "REQUIRED_DATE", "A Date is required")},
		))
		return -1, "", false
	}
	if !utils.ValidateDate(date) {
		utils.RespondJSON(w, http.StatusBadRequest, utils.CreateResponseError(
			"invalid_date",
			[]utils.APIError{utils.BuildApiError("date", "INVALID_DATE_FORMAT", "Date must be in YYYY-MM-DD format.")},
		))
		return -1, "", false
	}

	return userId, date, true
}

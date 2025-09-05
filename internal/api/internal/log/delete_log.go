package log

import (
	"net/http"
)

type DeleteLogEndPoint struct{}

func (DeleteLogEndPoint) GetLogId(w http.ResponseWriter, r *http.Request) *int {
	return GetLogId(w, r)
}

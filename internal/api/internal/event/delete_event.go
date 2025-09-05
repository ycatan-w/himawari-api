package event

import (
	"net/http"
)

type DeleteEventEndPoint struct{}

func (DeleteEventEndPoint) GetEventId(w http.ResponseWriter, r *http.Request) *int {
	return GetEventId(w, r)
}

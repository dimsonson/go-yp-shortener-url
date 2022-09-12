package handlers

import (
	"net/http"
)

func IncorrectRequests(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorect", http.StatusBadRequest)
}

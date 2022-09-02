package handlers

import (
	"net/http"
)

func DefHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "request incorect", http.StatusBadRequest)
}

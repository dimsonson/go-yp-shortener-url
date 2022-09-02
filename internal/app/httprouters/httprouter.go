package httprouters

import (
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
)

func HttpRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlers.GetHandler(w, r)
	case "POST":
		handlers.PostHandler(w, r)
	default:
		handlers.DefHandler(w, r)
	}
}

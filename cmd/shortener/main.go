package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {

	httprouters.HTTPRouter()
	log.Fatal(http.ListenAndServe(":8080", httprouters.NewRouter()))

}

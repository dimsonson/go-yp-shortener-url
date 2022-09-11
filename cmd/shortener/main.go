package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {
	addr := ":8080"
	httprouters.NewRouter()
	fmt.Printf("Started server on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, httprouters.NewRouter()))
}

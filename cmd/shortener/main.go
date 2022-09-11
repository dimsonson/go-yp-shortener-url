package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
)

func main() {

	httprouters.NewRouter()
	fmt.Println()
	log.Fatal(http.ListenAndServe(":8080", httprouters.NewRouter()))
	
}

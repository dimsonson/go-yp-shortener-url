package main

import (
	"log"
	"net/http"
)

// ShUrl — обработчик запроса.
func ShUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// если методом POST
	case "GET":
		// выдаем строку
		if err := r.ParseForm(); err != nil {
			// если не заполнена, возвращаем код ошибки
			http.Error(w, "Bad auth", 401)
			return
		}
	case "POST":
		// проверяем форму
		if err := r.ParseForm(); err != nil {
			// если не заполнена, возвращаем код ошибки
			http.Error(w, "Bad auth", 401)
			return
		}

	default:
		http.Error(w, "Bad auth", 401)
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", ShUrl)

	//server := &http.Server{
	//	Addr: "localhost:8080",
	//}

	// конструируем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
	//server.ListenAndServe()

}

/* type Middleware func(http.Handler) http.Handler

func Conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://yandex.ru/", http.StatusMovedPermanently)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// этот обработчик принимает только запросы, отправленные методом GET
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusMethodNotAllowed)
		return
	}
	// продолжаем обработку запроса
	// ...
} */

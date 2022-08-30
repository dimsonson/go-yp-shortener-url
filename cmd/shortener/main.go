package main

import (
	"log"
	"net/http"
)

// HelloWorld — обработчик запроса.
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="ru">
	  <head>
		<meta charset="utf-8">
		<title>Привет, мир!</title>
	  </head>
	  <body>
		<h1>Привет, мир!</h1>
		<p>Это веб-страница.</p>
	 </body>
	</html>`))
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HelloWorld)
	
	//server := &http.Server{
	//	Addr: "localhost:8080",
	//}

	// конструируем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
	//server.ListenAndServe()

}

type Middleware func(http.Handler) http.Handler

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
}

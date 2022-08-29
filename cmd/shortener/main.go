package main

import (
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
    // конструируем свой сервер
    server := &http.Server{
        Addr: "mydomain.com:80",
    }
    server.ListenAndServe()
} 
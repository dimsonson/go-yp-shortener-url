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
    // 
	//маршрутизация запросов обработчику
    http.HandleFunc("/", HelloWorld)
    // запуск сервера с адресом localhost, порт 8080
    http.ListenAndServe(":8080", nil)
}
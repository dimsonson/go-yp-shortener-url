package main

import (
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func main() {
	port := ":8080"
	log.Printf("Starting server on port %s\n", port)

	s := storage.NewMapStorage("map")
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs)
	r := httprouters.NewRouter(h)

	log.Fatal(http.ListenAndServe(port, r))
}

/* Задание для трека «Сервис сокращения URL»
Добавьте в сервер новый эндпоинт 
POST /api/shorten, принимающий 
в теле запроса JSON-объект
 {"url":"<some_url>"} 
 и возвращающий в ответ объект 
 {"result":"<shorten_url>"}.
Не забудьте добавить тесты на новый эндпоинт, 
как и на предыдущие.
Помните про HTTP content negotiation, 
проставляйте правильные значения в 
заголовок Content-Type.
{"url":"https://yandex.ru/search/?text=%D0%B7%D0%B0%D0%B3%D0%BE%D0%BB%D0%BE%D0%B2%D0%BE%D0%BA+http+%D0%BE%D1%82%D0%B2%D0%B5%D1%82%D0%B0+json&lr=213"} 

 */
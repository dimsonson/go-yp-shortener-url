package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func main() {
	// export SERVER_ADDRESS=localhost:8080
	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		log.Fatal("please, set SERVER_ADDRESS environment variable")
	}
	// export BASE_URL=localhost:8080
	_, ok = os.LookupEnv("BASE_URL")
	if !ok {
		log.Fatal("please, set BASE_URL environment variable")
	}

	log.Printf("starting server on %s\n", addr)

	s := storage.NewMapStorage("map")
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs)
	r := httprouters.NewRouter(h)

	log.Fatal(http.ListenAndServe(addr, r))
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

/*
 Инкремент 5
Задание для трека «Сервис сокращения URL»
Добавьте возможность конфигурировать сервис с помощью переменных окружения:
адрес запуска HTTP-сервера с помощью переменной SERVER_ADDRESS.
базовый адрес результирующего сокращённого URL с помощью переменной BASE_URL.
*/

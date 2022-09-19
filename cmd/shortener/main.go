package main

import (
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

func main() {
	// export SERVER_ADDRESS=localhost:8080
	// проверка переменной окуржения и присвоение значения по умолчанию, если не установлено
	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || !govalidator.IsURL(addr) {
		err := os.Setenv("SERVER_ADDRESS", "localhost:8080")
		if err != nil {
			log.Fatal("error setting default environment variable, please set SERVER_ADDRESS environment variable")
		}
		addr = os.Getenv("SERVER_ADDRESS")
	}
	log.Println("enviroment variable SERVER_ADDRESS set to defaulf value:", addr)

	// export BASE_URL=localhost:8080
	// проверка переменной окуржения и присвоение значения по умолчанию, если не установлено
	base, ok := os.LookupEnv("BASE_URL")
	if !ok || !govalidator.IsURL(base) {
		err := os.Setenv("BASE_URL", "http://localhost:8080")
		if err != nil {
			log.Fatal("error setting default environment variable, please set SERVER_ADDRESS environment variable")
		}
	}
	log.Println("enviroment variable BASE_URL set to defaulf value:", os.Getenv("BASE_URL"))

	// информирование, конфигурирование и запуск http сервера
	log.Printf("starting server on %s\n", addr)
	
	s := storage.NewFsStorage(make(map[string]string))
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

/* 
Задание для трека «Сервис сокращения URL»
Сохраняйте все сокращённые URL на диск в виде файла. При перезапуске приложения все URL должны быть восстановлены.
Путь до файла должен передаваться в переменной окружения FILE_STORAGE_PATH.
При отсутствии переменной окружения или при её пустом значении вернитесь к хранению сокращённых URL в памяти. 
*/
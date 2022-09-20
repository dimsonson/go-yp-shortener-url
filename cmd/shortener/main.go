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

	// проверка переменной окуржения и присвоение значения по умолчанию, если не установлено
	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || !govalidator.IsURL(addr) {
		err := os.Setenv("SERVER_ADDRESS", "localhost:8080")
		if err != nil {
			log.Fatal("error setting default environment variable, please set SERVER_ADDRESS environment variable")
		}
		addr = os.Getenv("SERVER_ADDRESS")
		log.Println("enviroment variable SERVER_ADDRESS set to default value:", addr)
	}

	// проверка переменной окуржения и присвоение значения по умолчанию, если не установлено
	base, ok := os.LookupEnv("BASE_URL")
	if !ok || !govalidator.IsURL(base) {
		err := os.Setenv("BASE_URL", "http://localhost:8080")
		if err != nil {
			log.Fatal("error setting default environment variable, please set SERVER_ADDRESS environment variable")
		}
		log.Println("enviroment variable BASE_URL set to default value:", os.Getenv("BASE_URL"))
	}

	// информирование, конфигурирование и запуск http сервера
	path, ok := os.LookupEnv("FILE_STORAGE_PATH")

	if !ok || !govalidator.IsUnixFilePath(path) {
		s := storage.NewMapStorage(make(map[string]string))
		log.Println("server will start with data storage in memory")
		srvs := services.NewService(s)
		h := handlers.NewHandler(srvs)
		r := httprouters.NewRouter(h)

		log.Printf("starting server on %s\n", addr)
		log.Fatal(http.ListenAndServe(addr, r))
		return
	}
	log.Println("server will start with data storage in file and memory cash")
	s := storage.NewFsStorage(make(map[string]string))
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs)
	r := httprouters.NewRouter(h)

	log.Printf("starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// export FILE_STORAGE_PATH=/db

// export BASE_URL=http://localhost:8080

// export SERVER_ADDRESS=localhost:8080

/* Инкремент 1
Задание для трека «Сервис сокращения URL»
Напишите сервис для сокращения длинных URL. Требования:
Сервер должен быть доступен по адресу: http://localhost:8080.
Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

/*
 Инкремент 2
Задание для трека «Сервис сокращения URL»
Покройте сервис юнит-тестами. Сконцентрируйтесь на покрытии тестами эндпоинтов, чтобы защитить API сервиса от случайных изменений.
*/

/*
 Инкремент 3
Задание для трека «Сервис сокращения URL»
Вы написали приложение с помощью стандартной библиотеки net/http. Используя любой пакет (роутер или фреймворк), совместимый с net/http, перепишите ваш код.
Задача направлена на рефакторинг приложения с помощью готовой библиотеки.
Обратите внимание, что необязательно запускать приложение вручную: тесты, которые вы написали до этого, помогут вам в рефакторинге.
*/

/*
 Инкремент 4
Задание для трека «Сервис сокращения URL»
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
 Инкремент 6
Задание для трека «Сервис сокращения URL»
Сохраняйте все сокращённые URL на диск в виде файла.
При перезапуске приложения все URL должны быть восстановлены.
Путь до файла должен передаваться в переменной окружения FILE_STORAGE_PATH.
При отсутствии переменной окружения или при её пустом значении вернитесь
к хранению сокращённых URL в памяти.
*/

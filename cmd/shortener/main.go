package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

const (
	defServAddr    = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
	defStoragePath = "db/keyvalue.json"
)

func main() {

	addrFlag := flag.String("a", defServAddr, "HTTP Server address")
	baseFlag := flag.String("b", defBaseURL, "Base URL")
	pathFlag := flag.String("f", defStoragePath, "File storage path")

	flag.Parse()

	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || !govalidator.IsURL(addr) {
		log.Println("eviroment variable SERVER_ADDRESS is empty or has wrong value ", addr)
		addr = *addrFlag
	}

	base, ok := os.LookupEnv("BASE_URL")
	if !ok || !govalidator.IsURL(base) {
		log.Println("eviroment variable BASE_URL is empty or has wrong value ", base)
		base = *baseFlag
	}

	path, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || !govalidator.IsUnixFilePath(path) || path == "" {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ", path)
		path = *pathFlag
	}

	var s services.StorageProvider

	if !govalidator.IsUnixFilePath(path) || path == "" {
		s = storage.NewMapStorage(make(map[string]string))
		log.Println("server will start with data storage in memory")

	} else {

		s = storage.NewFsStorage(make(map[string]string), path)
		log.Println("server will start with data storage in file and memory cash")

	}

	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs, base)
	r := httprouters.NewRouter(h)

	log.Printf("Base URL: %s\n", base)
	log.Printf("File storage path: %s\n", path)
	log.Printf("starting server on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// export FILE_STORAGE_PATH=db/keyvalue.json

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

/*
Инкремент 7
Задание для трека «Сервис сокращения URL»
Поддержите конфигурирование сервиса с помощью флагов командной строки наравне с уже имеющимися переменными окружения:
флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS);
флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL);
флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH).

*/

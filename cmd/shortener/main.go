package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var db = make(map[string]string)

// ShUrl — обработчик запроса
func ShUrl(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("обработчик")
	switch r.Method {
	// если метод POST
	case "GET":
		fmt.Println("get")
		fmt.Println(r.Method)
		fmt.Println(r.URL)
		fmt.Println(r.Header)

		u := strings.TrimPrefix(r.URL.Path, "/")
		fmt.Println(u)
		value, inMap := db[u]
		if !inMap {
			fmt.Println("нет такого URL") //
		}
		fmt.Println(value)

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		//w.Header().Add("Location", value)
		//устанавливаем заголовок Content-Type
		w.Header().Set("Location", value)
		fmt.Println(w.Header())
		// устанавливаем статус-код
		//http.Redirect(w, r, value, http.StatusTemporaryRedirect)
		w.WriteHeader(http.StatusTemporaryRedirect)
		fmt.Println(w.Header())
		// пишем тело ответа
		w.Write([]byte("test"))

	// если метод POST
	case "POST":
		//fmt.Println("post")
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		//fmt.Println(string(b))

		//создаем ключ
		key := randSeq(5)
		// проверяем уникальность ключа *********************

		//создаем пару ключ-значение
		db[key] = string(b)
		fmt.Println(db)

		//устанавливаем заголовок Content-Type
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		// устанавливаем статус-код 201
		//w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte("http://" + r.Host + "/" + key))
		fmt.Println(r.Host)

	default:
		http.Error(w, "Вы ввели неверный адрес", 401)
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

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
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
/*
var WebP1 = []byte(`<!DOCTYPE html>
<html>

<head>
  <meta charset="utf-8">
  <title>undefined</title>
  <meta name="generator" content="Google Web Designer 14.0.4.1108">
  <style type="text/css" id="gwd-text-style">
    p {
      margin: 0px;
    }
    h1 {
      margin: 0px;
    }
    h2 {
      margin: 0px;
    }
    h3 {
      margin: 0px;
    }
  </style>
  <style type="text/css">
    html, body {
      width: 100%;
      height: 100%;
      margin: 0px;
    }
    body {
      background-color: transparent;
      transform: perspective(1400px) matrix3d(1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1);
      transform-style: preserve-3d;
    }
    .gwd-input-1hwh {
      position: absolute;
      height: 1.79%;
      width: 30%;
      left: 34.53%;
      top: 16.09%;
    }
    .gwd-p-1izk {
      position: absolute;
      width: 384.375px;
      height: 37px;
      font-family: Roboto;
      font-weight: 500;
      top: 122px;
      left: 373px;
    }
  </style>
  <link href="https://fonts.googleapis.com/css?family=Roboto:100,100italic,300,300italic,regular,italic,500,500italic,700,700italic,900,900italic|Roboto+Mono:100,200,300,regular,500,600,700,100italic,200italic,300italic,italic,500italic,600italic,700italic|Source+Sans+Pro:200,200italic,300,300italic,regular,italic,600,600italic,700,700italic,900,900italic|Raleway:100,200,300,regular,500,600,700,800,900,100italic,200italic,300italic,italic,500italic,600italic,700italic,800italic,900italic|Noto+Sans:100,100italic,200,200italic,300,300italic,regular,italic,500,500italic,600,600italic,700,700italic,800,800italic,900,900italic|Mukta:200,300,regular,500,600,700,800|Inter:100,200,300,regular,500,600,700,800,900" rel="stylesheet" type="text/css">
</head>

<body class="htmlNoPages">
  <input type="text" id="text_1" class="gwd-input-1hwh">
  <p class="gwd-p-1izk">Вставьте ссылку для создания короткой версии:</p>
</body>

</html>`)
*/

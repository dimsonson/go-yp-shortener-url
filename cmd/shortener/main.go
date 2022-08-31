package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var db = make(map[string]string)
var keyLeght int = 5

// ShUrl — обработчик запроса
func ShortUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// если метод POST
	case "GET":
		// проверяем наличие ключа и получем длинную ссылку
		value, inMap := db[r.URL.Path]
		if !inMap {
			fmt.Println("нет такого URL")
		}
		http.Redirect(w, r, value, http.StatusTemporaryRedirect)
	// если метод POST
	case "POST":
		// читаем Body
		b, err := io.ReadAll(r.Body)
		// обрабатываем ошибку
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		//создаем ключ
		var key string //:= randSeq(5)
		// присваиваем значение ключа и проверяем уникальность ключа
		for {
			tmpKey := randSeq(keyLeght)
			if _, inMap := db[tmpKey]; inMap {
				return
			}
			key = tmpKey
			break
		}
		//создаем пару ключ-значение
		db[key] = string(b)
		//устанавливаем заголовок Content-Type
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		//устанавливаем статус-код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte("http://" + r.Host + key))
	default:
		http.Error(w, "Вы ввели неверный адрес", 400)
	}
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", ShortUrl)
	// конструируем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n+1)
	b[0] = rune('/')
	for i := range b[1:] {
		b[i+1] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

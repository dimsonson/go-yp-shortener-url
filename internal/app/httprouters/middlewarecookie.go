package httprouters

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/google/uuid"
)

// middleware функция распаковки-сжатия http алгоритмом gzip
func middlewareCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userid string
		userCookie, err := r.Cookie("token")
		// если токена нет в куке, токен не подписан, токена нет в хранилище - присвоение уникального userid
		if err != nil || userCookie.Value == "" || !TokenCheckSign(userCookie.Value, []byte(settings.SignKey)) {
			log.Println("request does not consist token cookie or empty - err:", err)
			userid = uuid.New().String()
			// подписание токена для возарата в ответе
			userTokenOut := TokenCreateSign(userid, []byte(settings.SignKey))
			// создаем куку
			cookie := &http.Cookie{
				Name:   "token",
				Value:  userTokenOut,
				MaxAge: 900,
			}
			// установим куку в ответ
			http.SetCookie(w, cookie)
		} else {
			// декодируем часть куки с userid
			useridByte, err := hex.DecodeString(userCookie.Value[:72])
			if err != nil {
				log.Printf("decodeString error: %v\n", err)
			}
			// приводим к string
			userid = string(useridByte)
		}
		// наследуем контекст, оснащаем его Value
		ctx := context.WithValue(r.Context(), settings.CtxKeyUserID , userid)
		// отправляем контекст дальше
		r = r.WithContext(ctx)
		// передаем запрос
		next.ServeHTTP(w, r)
	})
}

// проверка подписи iserid в куке
func TokenCheckSign(token string, key []byte) (ok bool) {
	// декодируем токен из строки в срез байтов
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		log.Printf("DecodeString error: %v\n", err)
	}
	//
	idBytes := tokenBytes[:36]

	h := hmac.New(sha256.New, key)
	h.Write(idBytes)

	newSign := h.Sum(nil)

	NewTokenBytes := append(idBytes, newSign[:]...)

	tokenNew := hex.EncodeToString(NewTokenBytes)

	ok = false
	if token == tokenNew {
		ok = true
	}

	log.Println("tokenCheckSign - ok :", ok)

	return ok
}

// создание куки с подписанным iserid
func TokenCreateSign(userid string, key []byte) (token string) {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(userid))

	sign := h.Sum(nil)

	tokenBytes := append([]byte(userid), sign[:]...)

	token = hex.EncodeToString(tokenBytes)

	return token
}

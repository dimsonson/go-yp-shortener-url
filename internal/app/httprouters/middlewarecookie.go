package httprouters

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/google/uuid"
)

// middleware функция распаковки-сжатия http алгоритмом gzip
func middlewareCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var userTokenIn string
		userCookie, err := r.Cookie("token")
		if err != nil {
			log.Println("MW _ Request does not consist token cookie - err:", err)
		} else {
			userTokenIn = userCookie.Value
		}

		// проверяем подпись токена
		var userid string
		if userTokenIn == "" {
			log.Println("userTokenIn is empty")
			userid = uuid.New().String()
		} else {
			userid, err = TokenCheckSign(userTokenIn, []byte(settings.SignKey))
			// если токена нет в куке, токен не подписан, токена нет в хранилище - присвоение уникального userid
			if err != nil {
				log.Println(err, "or userid doesnt exist in storage")
				userid = uuid.New().String()
			}
		}
		// наследуем контекcт запроса r *http.Request, оснащая его Timeout
		ctx := context.WithValue(r.Context(), "uuid", userid)
		//.WithTimeout(r.Context(), settings.StorageTimeout)
		// не забываем освободить ресурс
		r = r.WithContext(ctx)

		//defer cancel()
		//fmt.Print(ctx)

		// подписание токена для возарата в ответе
		userTokenOut := TokenCreateSign(userid, []byte(settings.SignKey))

		cookie := &http.Cookie{
			Name:   "token",
			Value:  userTokenOut,
			MaxAge: 300,
		}
		// установим куку в ответ
		http.SetCookie(w, cookie)
		next.ServeHTTP(w, r)

		/* // проверяем, что запрос содежит сжатые данные
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// читаем и распаковываем тело запроса с gzip
			gzR, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Println("gzip error: ", err)
			}
			r.Body = gzipReader{gzipReader: gzR, gzipBody: r.Body}
			defer gzR.Close()
			defer r.Body.Close()
		}
		// проверяем, что клиент поддерживает gzip-сжатие
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// создаём gzip.Writer поверх текущего w для записи сжатого ответа
			gzW, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				log.Println("gzip encodimg error:", err)
				return
			}
			defer gzW.Close()
			// устанавливаем заголовок сжатия содержимого ответа
			w.Header().Set("Content-Encoding", "gzip")
			// отправляем ответ с сжатым содержанием
			next.ServeHTTP(gzipWriter{ResponseWriter: w, gzWriter: gzW}, r)
			return
		}
		// если gzip не поддерживается клиентом, передаём управление дальше без изменений
		next.ServeHTTP(w, r) */
	})
}

// проверка подписи iserid в куке
func TokenCheckSign(token string, key []byte) (id string, err error) {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	fmt.Println("TokenCheckSign - token :", token)
	fmt.Println("TokenCheckSign - tokenBytes", tokenBytes)

	//id64 := binary.BigEndian.Uint64(tokenBytes[:36])

	//fmt.Println("TokenCheckSign - id64", id64)

	//idBytes := make([]byte, 36)
	idBytes := tokenBytes[:36]
	//	binary.BigEndian.PutUint64(idBytes, id64)

	fmt.Println("TokenCheckSign - idBytes", idBytes)

	h := hmac.New(sha256.New, key)
	h.Write(idBytes)

	fmt.Println("TokenCheckSign - h.Write(tokenBytes[:", tokenBytes)
	newSign := h.Sum(nil)
	fmt.Println("TokenCheckSign - newSign :", newSign)

	NewTokenBytes := append(idBytes, newSign[:]...)
	fmt.Println("TokenCheckSign - NewTokenBytes :", NewTokenBytes)

	tokenNew := hex.EncodeToString(NewTokenBytes)
	fmt.Println("TokenCheckSign - tokenNew :", tokenNew)
	if token != tokenNew {
		err = errors.New("sign incorrect")
	}
	id = fmt.Sprint(idBytes)
	fmt.Println("TokenCheckSign - err :", err)
	return id, err
}

// создание куки с подписанным iserid
func TokenCreateSign(userid string, key []byte) (token string) {
	fmt.Println("TokenCreateSign - userid :", userid)
	//uid := make([]byte, 36)
	//binary.BigEndian.PutUint32(uid, uint32(userid))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(userid))
	fmt.Println("TokenCreateSign - h.Write([]byte(uid)) :", []byte(userid))
	sign := h.Sum(nil)
	fmt.Println("TokenCreateSign - dst :", sign)

	tokenBytes := append([]byte(userid), sign[:]...)
	fmt.Println("TokenCreateSign - tokenBytes :", tokenBytes)

	token = hex.EncodeToString(tokenBytes)

	fmt.Println("TokenCreateSign - token :", token)

	return token
}

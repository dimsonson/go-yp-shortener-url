package httprouters

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
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
			ctx  := context.WithValue(r.Context(), "uuid", userid) 
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

// генерация и кодирование криптостойкого слайса байт
func RandomGenerator(n int) (cryproRand string, err error) {
	// определяем слайс нужной длины
	b := make([]byte, n)
	_, err = rand.Read(b) // записываем байты в массив b
	if err != nil {
		return ``, err
	}
	return hex.EncodeToString(b), nil
}

// проверка подписи iserid в куке
func TokenCheckSign(token string, key []byte) (id string, err error) {
	//tokenBytes := make([]byte, 5)
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	id32 := binary.BigEndian.Uint32(tokenBytes[:36])

	idBytes := make([]byte, 36)
	binary.BigEndian.PutUint32(idBytes, id32)

	h := hmac.New(sha256.New, key)
	h.Write(tokenBytes[:36])
	newSign := h.Sum(nil)

	NewTokenBytes := append(idBytes, newSign[:]...)
	tokenNew := hex.EncodeToString(NewTokenBytes)
	if token != tokenNew {
		err = errors.New("sign incorrect")
	}
	id = fmt.Sprint(id32)
	return id, err
}

// создание куки с подписанным iserid
func TokenCreateSign(userid string, key []byte) (token string) {

	uid := make([]byte, 36)
	//binary.BigEndian.PutUint32(uid, uint32(userid))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(uid))
	dst := h.Sum(nil)
	src := append(uid, dst[:]...)
	token = hex.EncodeToString(src)

	return token
}

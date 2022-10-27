package httprouters

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/google/uuid"
)

type ctxKey string
const keyUserID ctxKey = "uid"

// middleware функция распаковки-сжатия http алгоритмом gzip
func middlewareCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//var userTokenIn string
		var userid string
		userCookie, err := r.Cookie("token")
		// если токена нет в куке, токен не подписан, токена нет в хранилище - присвоение уникального userid
		if err != nil || userCookie.Value == "" || !TokenCheckSign(userCookie.Value, []byte(settings.SignKey)) {
			log.Println("request does not consist token cookie or empty - err:", err)
			userid = uuid.New().String()
			// подписание токена для возарата в ответе
			userTokenOut := TokenCreateSign(userid, []byte(settings.SignKey))
			//fmt.Println("middlewareCookie-userTokenOut :", userTokenOut)
			cookie := &http.Cookie{
				Name:   "token",
				Value:  userTokenOut,
				MaxAge: 900,
			}
			//fmt.Println("middlewareCookie-cookie :", cookie)
			// установим куку в ответ
			http.SetCookie(w, cookie)
		} else {
			// декодируем часть куки с userid
			useridByte, err := hex.DecodeString(userCookie.Value[:72])
			if err != nil {
				log.Printf("decodeString error: %v\n", err)
			}
			userid = string(useridByte)
			// fmt.Println("middlewareCookie-userid2 :", userid)
		}
		// наследуем контекст, оснащаем его Value
		ctx := context.WithValue(r.Context(), "uid", userid)
		fmt.Println("middlewareCookie-ctx :", ctx)
		// отправляем контекст дальше
		r = r.WithContext(ctx)
		// fmt.Println("middlewareCookie-userid3 :", userid)
		next.ServeHTTP(w, r)
	})
}

// проверка подписи iserid в куке
func TokenCheckSign(token string, key []byte) (ok bool) {
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		log.Printf("DecodeString error: %v\n", err)
	}

	//fmt.Println("TokenCheckSign - token :", token)
	//fmt.Println("TokenCheckSign - tokenBytes", tokenBytes)

	//id64 := binary.BigEndian.Uint64(tokenBytes[:36])

	//fmt.Println("TokenCheckSign - id64", id64)

	//idBytes := make([]byte, 36)
	idBytes := tokenBytes[:36]
	//	binary.BigEndian.PutUint64(idBytes, id64)

	//fmt.Println("TokenCheckSign - idBytes", idBytes)
	h := hmac.New(sha256.New, key)
	h.Write(idBytes)
	//fmt.Println("TokenCheckSign - h.Write(tokenBytes[:", tokenBytes)
	newSign := h.Sum(nil)
	//fmt.Println("TokenCheckSign - newSign :", newSign)
	NewTokenBytes := append(idBytes, newSign[:]...)
	//fmt.Println("TokenCheckSign - NewTokenBytes :", NewTokenBytes)
	tokenNew := hex.EncodeToString(NewTokenBytes)
	//fmt.Println("TokenCheckSign - tokenNew :", tokenNew)
	ok = false
	if token == tokenNew {
		//err = errors.New("sign incorrect")
		ok = true
	}
	//id = string(idBytes)
	fmt.Println("TokenCheckSign - ok :", ok)
	//fmt.Println("TokenCheckSign - id :", id)
	return ok
}

// создание куки с подписанным iserid
func TokenCreateSign(userid string, key []byte) (token string) {
	//fmt.Println("TokenCreateSign - userid :", userid)
	//uid := make([]byte, 36)
	//binary.BigEndian.PutUint32(uid, uint32(userid))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(userid))
	//fmt.Println("TokenCreateSign - h.Write([]byte(uid)) :", []byte(userid))
	sign := h.Sum(nil)
	//fmt.Println("TokenCreateSign - dst :", sign)
	tokenBytes := append([]byte(userid), sign[:]...)
	//fmt.Println("TokenCreateSign - tokenBytes :", tokenBytes)
	token = hex.EncodeToString(tokenBytes)
	//fmt.Println("TokenCreateSign - token :", token)
	return token
}

package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// интерфейс методов хранилища
type StorageProvider interface {
	PutToStorage(ctx context.Context, userid int, key string, value string) (err error)
	GetFromStorage(ctx context.Context, key string) (value string, err error)
	LenStorage(ctx context.Context) (lenn int)
	URLsByUserID(ctx context.Context, userid int) (userURLs map[string]string, err error)
	LoadFromFileToStorage()
	UserIDExist(ctx context.Context, userid int) bool
	StorageOkPing(ctx context.Context) (bool, error)
	StorageConnectionClose()
}

// структура конструктора бизнес логики
type Services struct {
	storage StorageProvider
}

// конструктор бизнес логики
func NewService(s StorageProvider) *Services {
	return &Services{
		s,
	}
}

// метод создание пары id : URL
func (sr *Services) ServiceCreateShortURL(ctx context.Context, url string, userTokenIn string) (key string, userTokenOut string) {
	// создаем и присваиваем значение короткой ссылки
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	var userid int
	if userTokenIn == "" {
		log.Println(err)
		userid = sr.storage.LenStorage(ctx)
	} else {
		userid, err = TokenCheckSign(userTokenIn, []byte(settings.SignKey))
		// если токена нет в куке, токен не подписан, токена нет в хранилище - присвоение уникального userid
		if err != nil || !sr.storage.UserIDExist(ctx, userid) {
			log.Println(err, "or userid doesnt exist in storage")
			userid = sr.storage.LenStorage(ctx)
		}
	}

	// подписание токена для возарата в ответе
	userTokenOut = TokenCreateSign(userid, []byte(settings.SignKey))

	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(ctx), key)
	// создаем пару ключ-значение в базе
	sr.storage.PutToStorage(ctx, userid, key, url)
	return key, userTokenOut
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(ctx context.Context, id string) (value string, err error) {
	// используем метод хранилища
	value, err = sr.storage.GetFromStorage(ctx, id)
	if err != nil {
		log.Println("id not found:", err)
	}
	return value, err
}

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(ctx context.Context,userToken string) (UserURLsMap map[string]string, err error) {
	//
	userid, err := TokenCheckSign(userToken, []byte(settings.SignKey))
	if err != nil {
		return UserURLsMap, err
	}
	//userid := userToken
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err := sr.storage.URLsByUserID(ctx, userid)
	if err != nil {
		log.Println(err)
		return map[string]string{"": ""}, err
	}
	fmt.Println("ServiceUserURLsMap:", userURLsMap)
	fmt.Println("ErrorServiceUserURLsMap:", err)

	return userURLsMap, err
}

// функция генерации случайной последовательности знаков
func RandSeq(n int) (string, error) {
	if n < 1 {
		err := fmt.Errorf("wromg argument: number %v less than 1\n ", n)
		return "", err
	}
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b), nil
}

// генерация и кодирование криптостойкого слайса байт
func RandomGenerator(n int) (string, error) {
	// определяем слайс нужной длины
	b := make([]byte, n)
	_, err := rand.Read(b) // записываем байты в массив b
	if err != nil {
		return ``, err
	}
	return hex.EncodeToString(b), nil
}

// проверка подписи iserid в куке
func TokenCheckSign(token string, key []byte) (userid int, err error) {
	//tokenBytes := make([]byte, 5)
	tokenBytes, err := hex.DecodeString(token)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("useridBytes", tokenBytes)

	id := binary.BigEndian.Uint32(tokenBytes[:4])
	fmt.Println("id :", id)

	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, id)

	h := hmac.New(sha256.New, key)
	h.Write(tokenBytes[:4])
	newSign := h.Sum(nil)
	fmt.Println("newSign", newSign)

	NewTokenBytes := append(idBytes, newSign[:]...)
	fmt.Printf("signed %x\n", newSign)
	fmt.Println("dst1append", NewTokenBytes)

	tokenNew := hex.EncodeToString(NewTokenBytes)
	fmt.Println("cookDst", token)

	if token != tokenNew {
		err = fmt.Errorf("sign incorrect")
	}

	return int(id), err
}

// создание куки с подписанным iserid
func TokenCreateSign(userid int, key []byte) (token string) {

	uid := make([]byte, 4)
	binary.BigEndian.PutUint32(uid, uint32(userid))
	fmt.Println("uid", uid)

	h := hmac.New(sha256.New, key)
	h.Write(uid)
	dst := h.Sum(nil)
	fmt.Println("dst", dst)

	src := append(uid, dst[:]...)
	fmt.Printf("signed %x\n", dst)
	fmt.Println("dst1append", src)

	token = hex.EncodeToString(src)
	fmt.Println("cookDst", token)

	return
}

func (sr *Services) ServiceStorageOkPing(ctx context.Context) (bool, error) {
	ok, err := sr.storage.StorageOkPing(ctx)
	return ok, err
}

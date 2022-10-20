package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// интерфейс методов хранилища
type StorageProvider interface {
	PutToStorage(ctx context.Context, userid int, key string, value string) (existKey string, err error)
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
func (sr *Services) ServiceCreateShortURL(ctx context.Context, url string, userTokenIn string) (key string, userTokenOut string, err error) {
	// создаем и присваиваем значение короткой ссылки
	key, err = RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	var userid int
	if userTokenIn == "" {
		log.Println("userTokenIn is empty")
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
	// создаем запись userid-ключ-значение в базе
	existKey, err := sr.storage.PutToStorage(ctx, userid, key, url)
	if err != nil {
		log.Println("request sr.storage.PutToStorage returned error:", err)
		key = existKey
	}
	return key, userTokenOut, err
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(ctx context.Context, id string) (value string, err error) {
	// используем метод хранилища
	value, err = sr.storage.GetFromStorage(ctx, id)
	if err != nil {
		log.Println("request sr.storage.GetFromStorageid returned error (id not found):", err)
	}
	return value, err
}

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(ctx context.Context, userToken string) (userURLsMap map[string]string, err error) {
	// проверяем подпись токена
	userid, err := TokenCheckSign(userToken, []byte(settings.SignKey))
	if err != nil {
		log.Println("token sign check returned error:", err)
		return
	}
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err = sr.storage.URLsByUserID(ctx, userid)
	if err != nil {
		log.Println("request sr.storage.URLsByUserID returned error:", err)
		return
	}
	return
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
	id := binary.BigEndian.Uint32(tokenBytes[:4])

	idBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(idBytes, id)

	h := hmac.New(sha256.New, key)
	h.Write(tokenBytes[:4])
	newSign := h.Sum(nil)

	NewTokenBytes := append(idBytes, newSign[:]...)
	tokenNew := hex.EncodeToString(NewTokenBytes)
	if token != tokenNew {
		err = errors.New("sign incorrect")
	}

	return int(id), err
}

// создание куки с подписанным iserid
func TokenCreateSign(userid int, key []byte) (token string) {

	uid := make([]byte, 4)
	binary.BigEndian.PutUint32(uid, uint32(userid))
	h := hmac.New(sha256.New, key)
	h.Write(uid)
	dst := h.Sum(nil)
	src := append(uid, dst[:]...)
	token = hex.EncodeToString(src)

	return
}

func (sr *Services) ServiceStorageOkPing(ctx context.Context) (bool, error) {
	ok, err := sr.storage.StorageOkPing(ctx)
	return ok, err
}

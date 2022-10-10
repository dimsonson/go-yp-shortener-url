package services

import (
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
	PutToStorage(userid string, key string, value string) (err error)
	GetFromStorage(key string) (value string, err error)
	LenStorage() (lenn int)
	URLsByUserID(userid string) (userURLs map[string]string, err error)
	LoadFromFileToStorage()
	UserIDExist(userid string) bool
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
func (sr *Services) ServiceCreateShortURL(url string, userTokenIn string) (key string, userTokenOut string) {
	// создаем и присваиваем значение короткой ссылки
	key, err := RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	// если токена нет в куке, токен не подписан, токена нет в хранилище - генерация и кодирование криптостойкого слайса байт
	if userTokenIn == "" || !TokenCheckSign(userTokenIn) || !sr.storage.UserIDExist(userTokenIn) {
		userTokenIn, err = RandomGenerator(settings.UserIDLeght)
		if err != nil {
			log.Println("error with userid generation : ", err)
		}
		fmt.Println("userTokenIn:", userTokenIn)
	}
	// подписание токена для возарата в ответе
	userTokenIn = userTokenOut
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(), key)
	// создаем пару ключ-значение в базе
	sr.storage.PutToStorage(userTokenIn, key, url)
	return key, userTokenOut
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(id string) (value string, err error) {
	// используем метод хранилища
	value, err = sr.storage.GetFromStorage(id)
	if err != nil {
		log.Println("id not found:", err)
	}
	return value, err
}

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(userToken string) (UserURLsMap map[string]string, err error) {

	//userid := userToken
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err := sr.storage.URLsByUserID(userToken)
	fmt.Println("userURLsMap:", userURLsMap)
	if err != nil {
		log.Println(err)
		return map[string]string{"": ""}, err
	}

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
func TokenCheckSign(userid int, token string, key []byte) bool {
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

	tokenNew := hex.EncodeToString(src)
	fmt.Println("cookDst", token)

	return tokenNew == token
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

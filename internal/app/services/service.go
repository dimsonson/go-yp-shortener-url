package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// интерфейс методов хранилища
type StorageProvider interface {
	PutToStorage(ctx context.Context, userid string, key string, value string) (string, error)
	GetFromStorage(ctx context.Context, key string) (string, error)
	LenStorage(ctx context.Context) int
	URLsByUserID(ctx context.Context, userid string) (map[string]string, error)
	LoadFromFileToStorage()
	UserIDExist(ctx context.Context, userid string) bool
	StorageOkPing(ctx context.Context) (bool, error)
	StorageConnectionClose()
	PutBatchToStorage(ctx context.Context, dc settings.DecodeBatchJSON) (dcCorr settings.DecodeBatchJSON, err error)
}

// структура конструктора бизнес логики
type Services struct {
	storage StorageProvider
	base    string
}

// конструктор бизнес логики
func NewService(s StorageProvider, base string) *Services {
	return &Services{
		s,
		base,
	}
}

// метод создание пары id : URL
func (sr *Services) ServiceCreateShortURL(ctx context.Context, url string) (key string, err error) {
	fmt.Println("ServiceCreateShortURL - url :::", url)
	// получаем значение из контекста
	//fmt.Println("ctx :", ctx.Value("uid").(string))
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	fmt.Println("ServiceCreateShortURL - ctx :", userid)

	// создаем и присваиваем значение короткой ссылки
	key, err = RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	fmt.Println("ServiceCreateShortURL key 1 ::: ", key)
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.LenStorage(ctx), key)
	fmt.Println("ServiceCreateShortURL key 2 ::: ", key)
	// создаем запись userid-ключ-значение в базе
	existKey, err := sr.storage.PutToStorage(ctx, userid, key, url)
	if err != nil {
		log.Println("request sr.storage.PutToStorage returned error:", err)
		key = existKey
	}
	fmt.Println("ServiceCreateShortURL key::: ", key)
	return key, err
}

// метод создание пакета пар id : URL
func (sr *Services) ServiceCreateBatchShortURLs(ctx context.Context, dc settings.DecodeBatchJSON) (ec []settings.EncodeBatch, err error) {
	// создаем и присваиваем значение короткой ссылки
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	fmt.Println(userid)

	// добавление shorturl
	for i := range dc {
		key, err := RandSeq(settings.KeyLeght)
		if err != nil {
			log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
		}
		key = fmt.Sprintf("%d%s", sr.storage.LenStorage(ctx), key)
		dc[i].ShortURL = key
	}
	fmt.Println("ServiceCreateBatchShortURLs dc", dc)

	//
	dc, err = sr.storage.PutBatchToStorage(ctx, dc)
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		break
	case err != nil:
		return nil, err
	}
	// заполняем слайс ответа EC := make([]settings.EncodeBatch,0)
	for _, v := range dc {
		elem := settings.EncodeBatch{
			CorrelationID: v.CorrelationID,
			ShortURL:      sr.base + "/" + v.ShortURL,
		}
		ec = append(ec, elem)
	}
	log.Println("HandlerCreateBatchJSON dc: ", dc)
	return ec, err
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(ctx context.Context, key string) (value string, err error) {
	// используем метод хранилища
	value, err = sr.storage.GetFromStorage(ctx, key)
	if err != nil {
		log.Println("request sr.storage.GetFromStorageid returned error (id not found):", err)
	}
	return value, err
}

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(ctx context.Context) (userURLsMap map[string]string, err error) {

	// получаем значение из контекста
	//fmt.Println("ctx :", ctx.Value("uid").(string))
	userid := ctx.Value(settings.CtxKeyUserID).(string)

	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err = sr.storage.URLsByUserID(ctx, userid)
	if err != nil {
		log.Println("request sr.storage.URLsByUserID returned error:", err)
		return userURLsMap, err
	}
	return userURLsMap, err
}

// функция генерации случайной последовательности знаков
func RandSeq(n int) (random string, ok error) {
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
	random = string(b)
	return random, nil
}

func (sr *Services) ServiceStorageOkPing(ctx context.Context) (ok bool, err error) {
	ok, err = sr.storage.StorageOkPing(ctx)
	return ok, err
}

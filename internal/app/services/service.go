package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// интерфейс методов хранилища
type StorageProvider interface {
	StoragePut(ctx context.Context, key string, value string, userid string) (existKey string, err error)
	StorageGet(ctx context.Context, key string) (value string, del bool, err error)
	StorageLen(ctx context.Context) (lenn int)
	StorageURLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error)
	StorageLoadFromFile()
	StorageOkPing(ctx context.Context) (bool, error)
	StorageConnectionClose()
	StoragePutBatch(ctx context.Context, dc settings.DecodeBatchJSON, userid string) (dcCorr settings.DecodeBatchJSON, err error)
	StorageDeleteURL(key string, userid string) (err error)
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
func (sr *Services) ServiceCreateShortURL(ctx context.Context, url string, userid string) (key string, err error) {
	// создаем и присваиваем значение короткой ссылки
	key, err = RandSeq(settings.KeyLeght)
	if err != nil {
		log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
	}
	// добавляем уникальный префикс к ключу
	key = fmt.Sprintf("%d%s", sr.storage.StorageLen(ctx), key)
	// создаем запись userid-ключ-значение в базе
	existKey, err := sr.storage.StoragePut(ctx, key, url, userid)
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		key = existKey
	case err != nil:
		return "", err
	}
	return key, err
}

// метод создание пакета пар id : URL
func (sr *Services) ServiceCreateBatchShortURLs(ctx context.Context, dc settings.DecodeBatchJSON, userid string) (ec []settings.EncodeBatch, err error) {
	// добавление shorturl
	for i := range dc {
		key, err := RandSeq(settings.KeyLeght)
		if err != nil {
			log.Fatal(err) //RandSeq настраивается на этапе запуска http сервера
		}
		key = fmt.Sprintf("%d%s", sr.storage.StorageLen(ctx), key)
		dc[i].ShortURL = key
	}
	// пишем в базу и получаем слайс с обновленными shorturl в случае конфликта
	dc, err = sr.storage.StoragePutBatch(ctx, dc, userid)
	switch {
	case err != nil && strings.Contains(err.Error(), "23505"):
		break
	case err != nil:
		return nil, err
	}
	// заполняем слайс ответа
	for _, v := range dc {
		elem := settings.EncodeBatch{
			CorrelationID: v.CorrelationID,
			ShortURL:      sr.base + "/" + v.ShortURL,
		}
		ec = append(ec, elem)
	}
	return ec, err
}

// метод возврат URL по id
func (sr *Services) ServiceGetShortURL(ctx context.Context, key string) (value string, del bool, err error) {
	// используем метод хранилища
	value, del, err = sr.storage.StorageGet(ctx, key)
	if err != nil {
		log.Println("request sr.storage.GetFromStorageid returned error (id not found):", err)
	}
	return value, del, err
}

// метод возврат всех URLs по userid
func (sr *Services) ServiceGetUserShortURLs(ctx context.Context, userid string) (userURLsMap map[string]string, err error) {
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err = sr.storage.StorageURLsByUserID(ctx, userid)
	if err != nil {
		log.Println("request sr.storage.URLsByUserID returned error:", err)
		return userURLsMap, err
	}
	return userURLsMap, err
}

// функция генерации псевдо случайной последовательности знаков
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

const workersCount = 1

// метод запись признака deleted_url
func (sr *Services) ServiceDeleteURL(shURLs [][2]string) {

	wg := &sync.WaitGroup{}
	inputCh := make(chan [2]string)
	// читаем массивы и кладём в inputCh
	wg.Add(1)
	go func() {
		for _, v := range shURLs {
			inputCh <- v
			fmt.Println("gorutine main to inputCh :", v)
		}
		wg.Done()
		close(inputCh)
	}()

	// здесь fanOut
	fanOutChs := fanOut(inputCh, workersCount)
	// var err error
	workerChs := make([]chan error, 0, workersCount)
	for _, fanOutCh := range fanOutChs {
		wg.Add(1)
		workerCh := make(chan error)

		// newWorker(fanOutCh, workerCh)
		// newWorker(input, out chan [2]string)
// нужен запуск нескольких воркеров, а не одного
		go func() {
			// обработка паники, что бы программа могла выполниться в этом случае
			/* 		defer func() {
				if x := recover(); x != nil {
					newWorker(fanOutCh, workerCh)
					log.Printf("run time panic: %v", x)
				}
			}() */

			for urls := range fanOutCh {
				err := sr.storage.StorageDeleteURL(urls[0], urls[1])

				workerCh <- err

				fmt.Println("worker out:", err)

			}

			close(workerCh)
		}()

		workerChs = append(workerChs, workerCh)
		wg.Done()
	}

	// здесь fanIn
	for v := range fanIn(workerChs...) {

	
			log.Println("delete request returned err: ", v)
	

	}
	wg.Wait()

}

func fanOut(inputCh chan [2]string, n int) []chan [2]string {
	chs := make([]chan [2]string, 0, n)
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		ch := make(chan [2]string)
		chs = append(chs, ch)
	}

	go func() {
		defer func(chs []chan [2]string) {
			for _, ch := range chs {
				close(ch)
			}
		}(chs)
		wg.Add(1)
		for i := 0; ; i++ {
			if i == len(chs) {
				i = 0
			}

			urls, ok := <-inputCh
			if !ok {
				return
			}

			ch := chs[i]
			ch <- urls
			fmt.Println("fanOut", urls)

		}
	}()
	wg.Wait()
	return chs
}

/* func newWorker(input, out chan [2]string) {
	go func() {
		// обработка паники, что бы программа могла выполниться в этом случае
		defer func() {
			if x := recover(); x != nil {
				newWorker(input, out)
				log.Printf("run time panic: %v", x)
			}
		}()

		for urls := range input {
			err := sr.storage.StorageDeleteURL(urls)

			out <- err

			fmt.Println("worker out:", err)

		}

		close(out)
	}()
}
 */
func fanIn(inputChs ...chan error) chan error {
	outCh := make(chan error)

	go func() {
		wg := &sync.WaitGroup{}

		for _, inputCh := range inputChs {
			wg.Add(1)

			go func(inputCh chan error) {
				defer wg.Done()
				for item := range inputCh {
					outCh <- item
					fmt.Println("fanIn to outCh :", item)

				}
			}(inputCh)
		}
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

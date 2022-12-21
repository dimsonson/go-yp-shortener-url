package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgerrcode"
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
	StoragePutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error)
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
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		key = existKey
	case err != nil:
		return "", err
	}
	return key, err
}

// метод создание пакета пар id : URL
func (sr *Services) ServiceCreateBatchShortURLs(ctx context.Context, dc models.BatchRequest, userid string) (ec []models.BatchResponse, err error) {
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
		elem := models.BatchResponse{
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

// метод запись признака deleted_url
func (sr *Services) ServiceDeleteURL(shURLs [][2]string) {
	// создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// создаем счетчик ожидания
	wg := &sync.WaitGroup{}
	// создаем выходной канал
	inputCh := make(chan [2]string)
	// горутина чтения массива и отправки ее значений в канал inputCh
	wg.Add(1)
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			log.Printf("stopped by cancel err : %v", ctx.Err())
			return
		default:
			for _, v := range shURLs {
				inputCh <- v
			}
			wg.Done()
			close(inputCh)
		}
	}(ctx)
	// здесь fanOut - получаем слайс каналов, в которые распределены значения из inputCh
	fanOutChs := fanOut(ctx, inputCh, settings.WorkersCount)
	// создаем слайс возвращемых каналов
	//workerChs := make([]chan error, 0, settings.WorkersCount)
	// итерируем по входным каналам с значениями и предаем из них значения в воркеры
	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan error)
		// запуск воркера
		wg.Add(1)
		go func(ctx context.Context, input chan [2]string, out chan error) {
			// итерация по входящим каналам воркера, выполнения обращения в хранилище
			for urls := range input {
				select {
				case <-ctx.Done():
					log.Printf("worker %s stopped by cancel err : %v", urls[0], ctx.Err())
					return
				default:
					err := sr.storage.StorageDeleteURL(urls[0], urls[1])
					// возвращаем значения из воркера в выходные каналы воркеров
					if err != nil {
						// логгируем в случае если из хранилища пришла ошибка
						log.Printf("shorturl %s from %s can't be delited with delete request error: %v", urls[0], urls[1], ctx.Err())
					}
					out <- err
				}
			}
			wg.Done()
			close(workerCh)
		}(ctx, fanOutCh, workerCh)
		// добавляем выходные каналы воркеров в слайс
	//	workerChs = append(workerChs, workerCh)
		//wg.Done()
	}
	// здесь fanIn - итерируем по слайсу каналов из воркеров и выводим их содержание в консоль
/* 	for v := range fanIn(ctx, workerChs...) {
		log.Println("delete request returned err: ", v)
	} */
	wg.Wait()
}

// функция распределения значений из одного канала в несколько по методу раунд робин
func fanOut(ctx context.Context, inputCh chan [2]string, n int) []chan [2]string {
	chs := make([]chan [2]string, 0, n)
	select {
	case <-ctx.Done():
		log.Printf("stopped by cancel err : %v", ctx.Err())
		return chs
	default:
		wg := &sync.WaitGroup{}
		for i := 0; i < n; i++ {
			ch := make(chan [2]string)
			chs = append(chs, ch)
		}
		go func(ctx context.Context) {
			defer func(chs []chan [2]string) {
				for _, ch := range chs {
					close(ch)
				}
			}(chs)
			wg.Add(1)
			select {
			case <-ctx.Done():
				log.Printf("stopped by cancel err : %v", ctx.Err())
				return
			default:
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
				}
			}
		}(ctx)
		wg.Wait()
		return chs
	}
}

/* // функция сбора значений из нескольких каналов в один
func fanIn(ctx context.Context, inputChs ...chan error) chan error {
	outCh := make(chan error)
	go func(ctx context.Context) {
		wg := &sync.WaitGroup{}
		select {
		case <-ctx.Done():
			log.Printf("stopped by cancel err : %v", ctx.Err())
			return
		default:
			for _, inputCh := range inputChs {
				wg.Add(1)
				go func(inputCh chan error) {
					defer wg.Done()
					for item := range inputCh {
						outCh <- item
					}
				}(inputCh)
			}
		}
		wg.Wait()
		close(outCh)
	}(ctx)
	return outCh
} */

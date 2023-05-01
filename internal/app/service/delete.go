// Package service пакет слоя бизнес логики.
package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
)

// DeleteStorageProvider интерфейс методов хранилища.
type DeleteStorageProvider interface {
	Delete(key string, userid string) (err error)
}

// DeleteServices структура конструктора бизнес логики.
type DeleteServices struct {
	storage DeleteStorageProvider
	base    string
}

// NewDeleteService конструктор бизнес  логики.
func NewDeleteService(s DeleteStorageProvider, base string) *DeleteServices {
	return &DeleteServices{
		s,
		base,
	}
}

// Delete метод записи признака deleted_url.
func (sr *DeleteServices) Delete(shURLs [][2]string) {
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
	// итерируем по входным каналам с значениями и предаем из них значения в воркеры
	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan error)
		// запуск воркера
		wg.Add(1)
		go func(ctx context.Context, input chan [2]string, out chan error) {
			// итерация по входящим каналам воркера, выполнения обращения в хранилище
			for urls := range input {
				// делаем паузу в соотвествии с Retry-After
				<-time.After(settings.RequestsTimeout)
				select {
				case <-ctx.Done():
					log.Printf("worker %s stopped by cancel err : %v", urls[0], ctx.Err())
					return
				default:
					err := sr.storage.Delete(urls[0], urls[1])
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
	}
	wg.Wait()
}

// fanOut функция распределения значений из одного канала в несколько по методу раунд робин.
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

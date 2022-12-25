package services

import (
	"context"
	"log"
)

// интерфейс методов хранилища
type GetStorageProvider interface {
	Get(ctx context.Context, key string) (value string, del bool, err error)
	GetBatch(ctx context.Context, userid string) (userURLs map[string]string, err error)
}

// структура конструктора бизнес логики
type GetServices struct {
	storage GetStorageProvider
	base    string
}

// конструктор бизнес  логики
func NewGetService(s GetStorageProvider, base string) *GetServices {
	return &GetServices{
		s,
		base,
	}
}

// метод возврат URL по id
func (sr *GetServices) Get(ctx context.Context, key string) (value string, del bool, err error) {
	// используем метод хранилища
	value, del, err = sr.storage.Get(ctx, key)
	if err != nil {
		log.Println("request sr.storage.GetFromStorageid returned error (id not found):", err)
	}
	return value, del, err
}

// метод возврат всех URLs по userid
func (sr *GetServices) GetBatch(ctx context.Context, userid string) (userURLsMap map[string]string, err error) {
	// используем метод хранилища для получения map URLs по userid
	userURLsMap, err = sr.storage.GetBatch(ctx, userid)
	if err != nil {
		log.Println("request sr.storage.URLsByUserID returned error:", err)
		return userURLsMap, err
	}
	return userURLsMap, err
}

package service

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
)

// PingStorageProvider интерфейс методов хранилища.
type PingStorageProvider interface {
	Ping(ctx context.Context) (bool, error)
	UsersQty(ctx context.Context) (usersQty int, err error)
	ShortsQty(ctx context.Context) (shortsQty int, err error)
}

// PingServices структура конструктора бизнес логики.
type PingServices struct {
	storage PingStorageProvider
}

// NewPingService конструктор бизнес  логики.
func NewPingService(s PingStorageProvider) *PingServices {
	return &PingServices{
		s,
	}
}

// Ping метод проверки достпности хранилища.
func (sr *PingServices) Ping(ctx context.Context) (ok bool, err error) {
	ok, err = sr.storage.Ping(ctx)
	return ok, err
}

// Stat метод получения статистики по пользователям и количеству обработанных ссылок.
func (sr *PingServices) Stat(ctx context.Context) (stat models.Stat, err error) {
	stat.Users, err = sr.storage.UsersQty(ctx)
	if err != nil {
		log.Print("sr.storage.UsersQty returned error::", err)
		return stat, err
	}
	stat.Urls, err = sr.storage.ShortsQty(ctx)
	if err != nil {
		log.Print("sr.storage.ShortsQty returned error::", err)
		return stat, err
	}
	return stat, err
}

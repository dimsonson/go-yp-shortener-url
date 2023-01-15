package service

import (
	"context"
)

// PingStorageProvider интерфейс методов хранилища.
type PingStorageProvider interface {
	Ping(ctx context.Context) (bool, error)
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

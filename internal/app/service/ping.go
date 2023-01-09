package services

import (
	"context"
)

// интерфейс методов хранилища
type PingStorageProvider interface {
	Ping(ctx context.Context) (bool, error)
}

// структура конструктора бизнес логики
type PingServices struct {
	storage PingStorageProvider
	base    string
}

// конструктор бизнес  логики
func NewPingService(s PingStorageProvider, base string) *PingServices {
	return &PingServices{
		s,
		base,
	}
}

func (sr *PingServices) Ping(ctx context.Context) (ok bool, err error) {
	ok, err = sr.storage.Ping(ctx)
	return ok, err
}
 
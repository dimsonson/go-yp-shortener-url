package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/rs/zerolog/log"
)

// PingServiceProvider интерфейс методов бизнес логики слоя Ping.
type PingServiceProvider interface {
	Ping(ctx context.Context) (bool, error)
	Stat(ctx context.Context) (stat models.Stat, err error)
}

// PingHandler структура для конструктура обработчика.
type PingHandler struct {
	service   PingServiceProvider
	trustCIDR *net.IPNet
}

// NewPingHandler конструктор обработчика.
func NewPingHandler(s PingServiceProvider, trustCIDR *net.IPNet) *PingHandler {
	return &PingHandler{
		s,
		trustCIDR,
	}
}

// Ping метод проверки доступности базы SQL.
func (hn PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	// создаем переменную для визуального возврата пользователю в теле отвта
	var result []byte
	ok, err := hn.service.Ping(ctx)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		result = []byte("DB ping NOT OK")
		log.Print(err)
	} else {
		w.WriteHeader(http.StatusOK)
		result = []byte("DB ping OK")
	}
	_, err = w.Write(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

// Stat метод обработки GET запроса из доверенной подсети и возвратом JSON c стат данными из хранилища.
func (hn PingHandler) Stat(w http.ResponseWriter, r *http.Request) {
	var err error
	if hn.trustCIDR != nil {
		// проверяем на соответсвие адреса отправителя в заголовке X-Real-IP доверенной подсети по CIDR
		ipFromHeader := r.Header.Get("X-Real-IP")
		reqIP := net.ParseIP(ipFromHeader)
		if !hn.trustCIDR.IP.Equal(reqIP.Mask(hn.trustCIDR.Mask)) {
			http.Error(w, "access denied", http.StatusForbidden)
			log.Print("X-Real-IP check error:", err)
			return
		}
	}
	// наследуем контекcт запроса r *http.Request, оснащая его Timeout
	ctx, cancel := context.WithTimeout(r.Context(), settings.StorageTimeout)
	// не забываем освободить ресурс
	defer cancel()
	//устанавливаем заголовок Content-Type
	w.Header().Set("content-type", "application/json; charset=utf-8")
	// запрос на получение stat - данных по количеству userid и short urls
	stat, err := hn.service.Stat(ctx)
	//устанавливаем статус-код 200 или 500
	switch {
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
	// пишем тело ответа
	err = json.NewEncoder(w).Encode(stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

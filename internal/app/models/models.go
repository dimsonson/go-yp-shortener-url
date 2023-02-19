// Package models пакет моделей структур используемых в различных пакетах.
package models

import "net"

// BatchRequest слайс структур декодирования JSON из POST запроса.
// type BatchRequest []struct {
// 	CorrelationID string `json:"correlation_id,omitempty"`
// 	OriginalURL   string `json:"original_url,omitempty"`
// 	ShortURL      string `json:"_"`
// }



type BatchRequest struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"_"`
}

// BatchResponse структура кодирования JSON для POST Batch ответа.
type BatchResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

// DecodeJSON структура декодирования JSON для POST запроса.
type DecodeJSON struct {
	URL string `json:"url,omitempty"`
}

// EncodeJSON структура кодирования JSON для POST запроса.
type EncodeJSON struct {
	Result string `json:"result,omitempty"`
}

// UserURL структура для создания среза surl:url и дельнейшего encode.
type UserURL struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}

// Config структура конфигурации сервиса, при запуске сервиса с флагом -c/config
// и отсутствии иных флагов и переменных окружения заполняется из файла указанного в этом флаге или переменной окружения CONFIG.
type Config struct {
	ServerAddress   string     `json:"server_address"`
	BaseURL         string     `json:"base_url"`
	FileStoragePath string     `json:"file_storage_path"`
	DatabaseDsn     string     `json:"database_dsn"`
	EnableHTTPS     bool       `json:"enable_https"`
	TrustedSubnet   string     `json:"trusted_subnet"`
	TrustedCIDR     *net.IPNet `json:"-"`
	ConfigJSON      string     `json:"-"`
}

// Stat структура для вывода статитстики по количеству сокращенных url и пользователей сервиса
type Stat struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

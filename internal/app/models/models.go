// Package models пакет моделей структур используемых в различных пакетах.
package models

// BatchRequest слайс структур декодирования JSON из POST запроса.
type BatchRequest []struct {
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

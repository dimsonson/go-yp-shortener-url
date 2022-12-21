package models

// слайс структур декодирования JSON из POST запроса
type BatchRequest []struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"_"`
}

// структура кодирования JSON для POST Batch ответа

type BatchResponse struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}
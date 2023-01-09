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

// структура декодирования JSON для POST запроса
type DecodeJSON struct {
	URL string `json:"url,omitempty"`
}

// структура кодирования JSON для POST запроса
type EncodeJSON struct {
	Result string `json:"result,omitempty"`
}


// структура для создания среза surl:url и дельнейшего encode
type UserURL struct {
	ShortURL    string `json:"short_url,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
}
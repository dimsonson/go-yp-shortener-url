package settings

import "time"

// длинна укороченной ссылки без первого слеш
const KeyLeght int = 5 //значение должно быть больше 0

// длинна слайса байт для UserID
const UserIDLeght int = 16 //значение должно быть больше 0

// ключ подписи
const SignKey string = "9e9e0b4e6de418b2f84fca35165571c5"

// timeout запроса
const StorageTimeout = 300000 * time.Second

// имя таблицы в базе PosgreSQL
const SQLTableName = "sh_urls"

// тип для context.WithValue
type ctxKey string
// ключ для context.WithValue
const CtxKeyUserID ctxKey = "uid"

// слайс структур декодирования JSON из POST запроса
type DecodeBatchJSON []struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"_"`
}

// структура кодирования JSON для POST Batch ответа

type EncodeBatch struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
} 


// количество каналов для воркеров при установке пометку удаленный для sh_urls 
const WorkersCount = 30

// константы цветового вывода в консоль
const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)
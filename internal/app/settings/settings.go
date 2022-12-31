// settings пакет настроек сервиса.
package settings

import "time"

// RequestsTimeout начальный таймаут для горутин.
var RequestsTimeout = 100 * time.Millisecond

// KeyLeght длинна укороченной ссылки без первого слеш.
const KeyLeght int = 5 //значение должно быть больше 0

// UserIDLeght длинна слайса байт для UserID.
const UserIDLeght int = 16 //значение должно быть больше 0

// SignKey ключ подписи.
const SignKey string = "9e9e0b4e6de418b2f84fca35165571c5"

// StorageTimeout timeout запроса.
const StorageTimeout = 5 * time.Second

// SQLTableName имя таблицы в базе PosgreSQL.
const SQLTableName = "sh_urls"

// ctxKey тип для context.WithValue.
type ctxKey string

// ctxKey ключ для context.WithValue.
const CtxKeyUserID ctxKey = "uid"

// WorkersCount количество каналов для воркеров при установке пометку удаленный для sh_urls.
const WorkersCount = 30

// WorkersCount константы цветового вывода в консоль.
const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

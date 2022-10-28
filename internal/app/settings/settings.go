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

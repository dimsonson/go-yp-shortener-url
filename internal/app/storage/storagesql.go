package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// структура хранилища
type StorageSQL struct {
	PostgreSQL *sql.DB
	//buffer     settings.DecodeBatchJSON // для 4 спринта
}

// метод записи id:url в хранилище
func (ms *StorageSQL) PutToStorage(ctx context.Context, key string, value string) (existKey string, err error) {
	// получаем значение iserid из контекста
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	// столбец short_url в SQL таблице содержит только уникальные занчения
	// создаем текст запроса
	q := `INSERT INTO sh_urls 
			VALUES (
			$1, 
			$2, 
			$3
			)`
	// записываем в хранилице userid, id, URL
	_, err = ms.PostgreSQL.ExecContext(ctx, q, userid, key, value)
	if err != nil {
		log.Println("insert SQL request PutToStorage scan error:", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				// создаем текст запроса
				q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
				// запрос в хранилище на корокий URL по длинному URL,
				// пишем результат запроса в пременную existKey
				err := ms.PostgreSQL.QueryRowContext(ctx, q, value).Scan(&existKey)
				if err != nil {
					log.Println("select PutToStorage SQL select request scan error:", err)
				}
			}
		}
	}
	return existKey, err
}

// метод пакетной записи id:url в хранилище
func (ms *StorageSQL) PutBatchToStorage(ctx context.Context, dc settings.DecodeBatchJSON) (dcCorr settings.DecodeBatchJSON, err error) {
	// получаем значение iserid из контекста
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	// готовим инструкцию
	q := "INSERT INTO sh_urls VALUES ($1, $2, $3)"
	// итерируем по слайсу структур
	for i, v := range dc {
		// добавляем значения в транзакцию
		_, err = ms.PostgreSQL.ExecContext(ctx, q, userid, v.ShortURL, v.OriginalURL)
		if err != nil {
			log.Println("insert SQL request PutBatchToStorage scan error:", err)
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case pgerrcode.UniqueViolation:
					fmt.Println("pgErr.Code:::", pgErr.Code)
					// создаем текст запроса
					q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
					// запрос в хранилище на корокий URL по длинному URL,
					// пишем результат запроса в пременную existKey
					err = ms.PostgreSQL.QueryRowContext(ctx, q, v.OriginalURL).Scan(&dc[i].ShortURL)
					if err != nil {
						log.Println("PutBatchToStorage select SQL request scan error:", err)
					}
				}
			}
		}
	}

	return dc, err
}

// конструктор нового хранилища PostgreSQL
func NewSQLStorage(p string) *StorageSQL {
	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()
	// открываем базу данных
	db, err := sql.Open("pgx", p)
	if err != nil {
		log.Println("database opening error:", err)
	}
	// создаем текст запроса
	// возможно ли имя таблицы вывести в файл settings?
	q := `CREATE TABLE IF NOT EXISTS sh_urls (
				"userid" TEXT,
				"short_url" TEXT NOT NULL UNIQUE,
				"long_url" TEXT NOT NULL UNIQUE
				)`
	// создаем таблицу в SQL базе, если не существует
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		log.Println("request NewSQLStorage to sql db returned error:", err)
	}
	return &StorageSQL{
		PostgreSQL: db,
		// buffer: make(settings.DecodeBatchJSON, 0, settings.BufferBatchSQL),  // для 4 спринта
	}
}

// метод получения записи из хранилища
func (ms *StorageSQL) GetFromStorage(ctx context.Context, key string) (value string, err error) {
	// создаем текст запроса
	q := `SELECT long_url FROM sh_urls WHERE short_url = $1`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRowContext(ctx, q, key)
	// пишем результат запроса в пременную value
	err = row.Scan(&value)
	if err != nil {
		log.Println("select GetFromStorage SQL request scan error:", err)
		return value, err
	}
	return value, err
}

// метод определения длинны хранилища
func (ms *StorageSQL) LenStorage(ctx context.Context) (lenn int) {
	// создаем текст запроса
	q := `SELECT COUNT (*) FROM sh_urls`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRowContext(ctx, q)
	// пишем результат запроса в пременную lenn
	err := row.Scan(&lenn)
	if err != nil {
		log.Println("select LenStorage SQL request scan error:", err)
	}
	return lenn
}

// метод отбора URLs по UserID
// посмотреть возможность использования SQLx
func (ms *StorageSQL) URLsByUserID(ctx context.Context) (userURLs map[string]string, err error) {
	// получаем значение iserid из контекста
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	// создаем текст запроса
	q := `SELECT short_url, long_url FROM sh_urls WHERE userid = $1`
	// делаем запрос в SQL, получаем строку
	rows, err := ms.PostgreSQL.QueryContext(ctx, q, userid)
	if err != nil {
		log.Println("select URLsByUserID SQL reuest error :", err)
	}
	defer rows.Close()
	// пишем результат запроса в map
	userURLs = make(map[string]string)
	for rows.Next() {
		var k, v string
		err = rows.Scan(&k, &v)
		if err != nil {
			log.Println("row by row scan URLsByUserID error :", err)
		}
		userURLs[k] = v
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Println("request URLsByUserID iteration scan error:", err)
	}
	// проверяем наличие записей
	if len(userURLs) == 0 {
		err = fmt.Errorf("userid not found in the storage")
	}
	return userURLs, err
}

func (ms *StorageSQL) LoadFromFileToStorage() {
}

// пинг хранилища для api/user/urls
func (ms *StorageSQL) StorageOkPing(ctx context.Context) (ok bool, err error) {
	err = ms.PostgreSQL.PingContext(ctx)
	if err != nil {
		return false, err
	}
	return true, err
}

// метод закрытия совединения с SQL базой
func (ms *StorageSQL) StorageConnectionClose() {
	ms.PostgreSQL.Close()
}

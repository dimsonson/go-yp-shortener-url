package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// структура хранилища
type StorageSQL struct {
	PostgreSQL *sql.DB
}

// метод записи id:url в хранилище
func (ms *StorageSQL) StoragePut(ctx context.Context, key string, value string, userid string) (existKey string, err error) {

	// создаем текст запроса
	q := `INSERT INTO sh_urls 
			VALUES (
			$1, 
			$2, 
			$3,
			$4
			)`
	// записываем в хранилице userid, id, URL
	_, err = ms.PostgreSQL.ExecContext(ctx, q, userid, key, value, false)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// создаем текст запроса
		q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
		// запрос в хранилище на корокий URL по длинному URL, пишем результат запроса в пременную existKey
		err := ms.PostgreSQL.QueryRowContext(ctx, q, value).Scan(&existKey)
		if err != nil {
			log.Println("select PutToStorage SQL select request scan error:", err)
			return "", err
		}
	}
	if err != nil {
		log.Println("insert SQL request PutToStorage scan error:", err)
		return "", err
	}
	return existKey, err
}

// метод пакетной записи id:url в хранилище
func (ms *StorageSQL) StoragePutBatch(ctx context.Context, dc models.BatchRequest, userid string) (dcCorr models.BatchRequest, err error) {
	// объявляем транзакцию
	tx, err := ms.PostgreSQL.Begin()
	if err != nil {
		log.Println("error PutBatchToStorage tx.Begin : ", err)
		return nil, err
	}
	// если возникает ошибка, откатываем изменения
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Println("insert urls: unable to rollback: ", rollbackErr)
		return nil, err
	}
	// готовим инструкцию
	stmt, err := ms.PostgreSQL.PrepareContext(ctx, "INSERT INTO sh_urls VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Println("error PutBatchToStorage stmt.PrepareContext : ", err)
		return nil, err
	}
	// закрываем инструкцию, когда она больше не нужна
	defer stmt.Close()
	// итерируем по слайсу структур
	var pgErr *pgconn.PgError
	for i, v := range dc {
		// добавляем значения в транзакцию
		_, err = stmt.ExecContext(ctx, userid, v.ShortURL, v.OriginalURL, false)
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			// создаем текст запроса
			q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
			// запрос в хранилище на корокий URL по длинному URL, пишем результат запроса в поле структуры
			err := ms.PostgreSQL.QueryRowContext(ctx, q, v.OriginalURL).Scan(&dc[i].ShortURL)
			if err != nil {
				log.Println("select SQL request PutBatchToStorage scan error:", err)
				return nil, err
			}
		}
		if err != nil {
			log.Println("insert SQL request PutBatchToStorage scan error:", err)
			return nil, err
		}
	}
	// сохраняем изменения
	if err = tx.Commit(); err != nil {
		log.Println("error PutBatchToStorage tx.Commit : ", err)
	}
	return dcCorr, err
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
		log.Println("database opening error:", settings.ColorRed, err, settings.ColorReset)
	}
	// создаем текст запроса
	q := `CREATE TABLE IF NOT EXISTS sh_urls (
				"userid" TEXT,
				"short_url" TEXT NOT NULL UNIQUE,
				"long_url" TEXT NOT NULL UNIQUE,
				"deleted_url" BOOLEAN 
				)`
	// создаем таблицу в SQL базе, если не существует
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		log.Println("request NewSQLStorage to sql db returned error:", settings.ColorRed, err, settings.ColorReset)
	}
	return &StorageSQL{
		PostgreSQL: db,
	}
}

// метод получения записи из хранилища
func (ms *StorageSQL) StorageGet(ctx context.Context, key string) (value string, del bool, err error) {
	// создаем текст запроса
	q := `SELECT long_url, deleted_url FROM sh_urls WHERE short_url = $1`
	// делаем запрос в SQL, получаем строку и пишем результат запроса в пременную value
	err = ms.PostgreSQL.QueryRowContext(ctx, q, key).Scan(&value, &del)
	if err != nil {
		log.Println("select GetFromStorage SQL request scan error:", err)
		return value, del, err
	}
	return value, del, err
}

// метод определения длинны хранилища
func (ms *StorageSQL) StorageLen(ctx context.Context) (lenn int) {
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
func (ms *StorageSQL) StorageURLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {
	// получаем значение iserid из контекста
	// userid := ctx.Value(settings.CtxKeyUserID).(string)
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

func (ms *StorageSQL) StorageLoadFromFile() {
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

// метод запись признака deleted_url
func (ms *StorageSQL) StorageDeleteURL(key string, userid string) (err error) {
	q := `UPDATE sh_urls SET deleted_url = true WHERE short_url = $1 AND userid = $2`
	// записываем в хранилице userid, id, URL
	_, err = ms.PostgreSQL.Exec(q, key, userid)
	if err != nil {
		log.Println("update SQL request StorageDeleteURL error:", err)
	}
	return err
}

package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// структура хранилища
type StorageSQL struct {
	PostgreSQL *pgxpool.Pool
	buffer     settings.DecodeBatchJSON
}

// метод записи id:url в хранилище
func (ms *StorageSQL) PutToStorage(ctx context.Context, userid string, key string, value string) (existKey string, err error) {
	fmt.Println("PutToStorage userid, key, value:::", userid, key, value)

	// столбец short_url в SQL таблице содержит только уникальные занчения
	// создаем текст запроса
	q := `INSERT INTO sh_urls 
			VALUES (
			$1, 
			$2, 
			$3
			)`
	// записываем в хранилице userid, id, URL
	_, err = ms.PostgreSQL.Exec(ctx, q, userid, key, value)
	if err != nil {
		log.Println("insert request PutToStorage scan error:", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				log.Println("correctly matched ", pgErr.Code)
				// создаем текст запроса
				q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
				// запрос в хранилище на корокий URL по длинному URL,
				// пишем результат запроса в пременную existKey
				err := ms.PostgreSQL.QueryRow(ctx, q, value).Scan(&existKey)
				if err != nil {
					log.Println("PutToStorage select request scan error:", err)
				}
			}
		}
	}
	fmt.Println("PutToStorage existKey::: ", existKey)
	return existKey, err
}

// метод пакетной записи id:url в хранилище
func (ms *StorageSQL) PutBatchToStorage(ctx context.Context, dc settings.DecodeBatchJSON) (dcCorr settings.DecodeBatchJSON, err error) {
	userid := ctx.Value(settings.CtxKeyUserID).(string)
	fmt.Println(userid)
	fmt.Println("dc", dc)
	// готовим инструкцию
	q := "INSERT INTO sh_urls VALUES ($1, $2, $3)"
	// итерируем по слайсу структур
	for i, v := range dc {
		// добавляем значения в транзакцию
		if _, err = ms.PostgreSQL.Exec(ctx, q, userid, v.ShortURL, v.OriginalURL); err != nil {
			if err != nil {
				log.Println("insert request PutToStorage scan error:", err)
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) {
					switch pgErr.Code {
					case pgerrcode.UniqueViolation:
						log.Println("correctly matched ", pgErr.Code)
						// создаем текст запроса
						q := `SELECT short_url FROM sh_urls WHERE long_url = $1`
						// запрос в хранилище на корокий URL по длинному URL,
						// пишем результат запроса в пременную existKey
						err := ms.PostgreSQL.QueryRow(ctx, q, v.OriginalURL).Scan(&dc[i].ShortURL)
						if err != nil {
							log.Println("PutToStorage select request scan error:", err)
						}
					}
				}
			}
		}
	}

	//fmt.Println("ms.buffer", ms.buffer)
	fmt.Println("Storage SQL dc", dc)
	return dc, err
}

// конструктор нового хранилища PostgreSQL
func NewSQLStorage(p string) *StorageSQL {
	// создаем контекст и оснащаем его таймаутом
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()
	// открываем базу данных
	dbpool, err := pgxpool.Connect(ctx, p)
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
	_, err = dbpool.Exec(ctx, q)
	if err != nil {
		log.Println("request NewSQLStorage to sql db returned error:", err)
	}
	return &StorageSQL{
		PostgreSQL: dbpool,
		buffer:     make(settings.DecodeBatchJSON, 0, settings.BufferBatchSQL),
	}
}

// метод получения записи из хранилища
func (ms *StorageSQL) GetFromStorage(ctx context.Context, key string) (value string, err error) {
	// создаем текст запроса
	q := `SELECT long_url FROM sh_urls WHERE short_url = $1`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRow(ctx, q, key)
	// пишем результат запроса в пременную value
	err = row.Scan(&value)
	if err != nil {
		log.Println("scan GetFromStorage to value variable returned error:", err)
		return value, err
	}
	return value, err
}

// метод определения длинны хранилища
func (ms *StorageSQL) LenStorage(ctx context.Context) (lenn int) {
	// создаем текст запроса
	q := `SELECT COUNT (*) FROM sh_urls`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRow(ctx, q)
	// пишем результат запроса в пременную lenn
	err := row.Scan(&lenn)
	if err != nil {
		log.Println("scan LenStorage to lenn variable returned error:", err)
	}
	return lenn
}

// метод отбора URLs по UserID
// посмотреть возможность использования SQLx
func (ms *StorageSQL) URLsByUserID(ctx context.Context, userid string) (userURLs map[string]string, err error) {
	// создаем текст запроса
	q := `SELECT short_url, long_url FROM sh_urls WHERE userid = $1`
	// делаем запрос в SQL, получаем строку
	rows, err := ms.PostgreSQL.Query(ctx, q, userid)
	if err != nil {
		log.Println("sql reuest URLsByUserID error :", err)
	}
	defer rows.Close()
	// пишем результат запроса в map
	userURLs = make(map[string]string)
	for rows.Next() {
		var k, v string
		err = rows.Scan(&k, &v)
		if err != nil {
			log.Println("row scan URLsByUserID error :", err)
		}
		userURLs[k] = v
	}
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Println("request URLsByUserID iteration scan error:", err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		err := errors.New("request URLsByUserID has no content for this token")
		return nil, err
	}
	return userURLs, err
}

func (ms *StorageSQL) LoadFromFileToStorage() {
}

// посик userid в хранилице
func (ms *StorageSQL) UserIDExist(ctx context.Context, userid string) bool {
	var DBuserid string
	q := `SELECT userid from sh_urls WHERE userid = $1`
	row := ms.PostgreSQL.QueryRow(ctx, q, userid)
	err := row.Scan(&DBuserid)
	if err != nil {
		log.Println("request UserIDExist returned scan error:", err)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false
	}
	return true
}

// пинг хранилища для api/user/urls
func (ms *StorageSQL) StorageOkPing(ctx context.Context) (ok bool, err error) {
	err = ms.PostgreSQL.Ping(ctx)
	if err != nil {
		return false, err
	}
	return true, err
}

// метод закрытия совединения с SQL базой
func (ms *StorageSQL) StorageConnectionClose() {
	ms.PostgreSQL.Close()
}

// метод пакетной записи в базу из буфера
/* func (ms *StorageSQL) Flush(ctx context.Context, userid string) error {
	// проверим на всякий случай
	if ms.PostgreSQL == nil {
		return errors.New("you haven`t opened the database connection")
	}
	tx, err := ms.PostgreSQL.BeginTx(ctx)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO videos(title, description, views, likes) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range ms.buffer {
		if _, err = stmt.Exec(userid, v.OriginalURL); err != nil {
			if err = tx.Rollback(); err != nil {
				log.Fatalf("update drivers: unable to rollback: %v", err)
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("update drivers: unable to commit: %v", err)
		return err
	}

	ms.buffer = ms.buffer[:0]
	return nil
}
*/

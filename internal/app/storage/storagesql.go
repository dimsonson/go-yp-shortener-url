package storage

import (
	"context"
	"errors"
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
}

// метод записи id:url в хранилище
func (ms *StorageSQL) PutToStorage(ctx context.Context, userid int, key string, value string) (existKey string, err error) {
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
	return existKey, err
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
				"userid" INTEGER,
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
func (ms *StorageSQL) URLsByUserID(ctx context.Context, userid int) (userURLs map[string]string, err error) {
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
func (ms *StorageSQL) UserIDExist(ctx context.Context, userid int) bool {
	var DBuserid int
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

// сервис закрытия совединения с SQL базой
func (ms *StorageSQL) StorageConnectionClose() {
	ms.PostgreSQL.Close()
}

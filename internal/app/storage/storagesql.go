package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// структура хранилища
type StorageSQL struct {
	PostgreSQL *sql.DB
}

// метод записи id:url в хранилище
func (ms *StorageSQL) PutToStorage(ctx context.Context, userid int, key string, value string) (err error) {
	// столбец short_url в SQL таблице содержит только иниткальные занчения
	// создаем текст запроса
	q := `INSERT INTO sh_urls 
			VALUES (
			$1, 
			$2, 
			$3
			)`

	// записываем в хранилице userid, id, URL
	res, err := ms.PostgreSQL.ExecContext(ctx, q, userid, key, value)
	if err != nil {
		return err
	}

	log.Println("PutToStorage: ", res)

	return nil
}

// конструктор нового хранилища JSON
func NewSQLStorage(p string) (*StorageSQL, *sql.DB) {
	db, err := sql.Open("pgx", p)
	if err != nil {
		log.Println(err)
	}
	//defer db.Close()
	// создание таблицы SQL если не существует
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()
	// создаем текст запроса
	q := `CREATE TABLE IF NOT EXISTS sh_urls (
				"userid" INTEGER,
				"short_url" TEXT NOT NULL UNIQUE,
				"long_url" TEXT NOT NULL UNIQUE
			)`

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		log.Println(err)
	}

	return &StorageSQL{
		PostgreSQL: db,
	}, db
}

// метод получения записи из хранилища
func (ms *StorageSQL) GetFromStorage(ctx context.Context, key string) (value string, err error) {
	// создаем текст запроса
	q := `SELECT long_url FROM sh_urls WHERE short_url = $1`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRowContext(ctx, q, key)
	// пишем результат запроса в пременную lenn
	err = row.Scan(&value)
	if err != nil {
		log.Println("SQL request scan error:", err)
		return "", err
	}

	fmt.Println("value:", value)

	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageSQL) LenStorage(ctx context.Context) (lenn int) {
	// создаем текст запроса
	q := `SELECT COUNT(*) FROM sh_urls`
	// делаем запрос в SQL, получаем строку
	row := ms.PostgreSQL.QueryRowContext(ctx, q)
	// пишем результат запроса в пременную lenn
	err := row.Scan(&lenn)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Lenn:", lenn)
	return lenn
}

// метод отбора URLs по UserID
// посмотреть возможность использования SQLx
func (ms *StorageSQL) URLsByUserID(ctx context.Context, userid int) (userURLs map[string]string, err error) {
	// создаем текст запроса
	q := `SELECT short_url, long_url FROM sh_urls WHERE userid = $1`
	// делаем запрос в SQL, получаем строку
	rows, err := ms.PostgreSQL.QueryContext(ctx, q, userid)
	if err != nil {
		log.Println("sql reuest URLsByUserID error :", err)
	}
	defer rows.Close()
	fmt.Println("1ErrorURLsByUserIDService:", err)
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
	fmt.Println("2ErrorURLsByUserIDService:", err)
	// проверяем итерации на ошибки
	err = rows.Err()
	if err != nil {
		log.Println("SQL request scan error:", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		err := fmt.Errorf("no content for this token")
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
	row := ms.PostgreSQL.QueryRowContext(ctx, q, userid)
	err := row.Scan(&DBuserid)
	if err != nil {
		log.Println("SQL request scan error:", err)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	return true
}

// пинг хранилища для api/user/urls
func (ms *StorageSQL) StorageOkPing(ctx context.Context) (bool, error) {
	err := ms.PostgreSQL.PingContext(ctx)
	if err != nil {
		return false, err
	}
	return true, err
}

// сервис закрытия совединения с SQL базой
func (ms *StorageSQL) StorageConnectionClose() {
	ms.PostgreSQL.Close()
}

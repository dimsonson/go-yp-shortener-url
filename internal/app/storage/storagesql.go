package storage

import (
	"context"
	"database/sql"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// структура хранилища
type StorageSQL struct {
	//UserID   map[string]int    `json:"iserid,omitempty"` // shorturl:userid
	//IDURL    map[string]string `json:"idurl,omitempty"`  // shorturl:URL
	//pathName string
	PostgreSQL string //*sql.DB
}

// метод записи id:url в хранилище
func (ms *StorageSQL) PutToStorage(userid int, key string, value string) (err error) {
	/* // проверяем наличие ключа в хранилище
	if value, ok := ms.IDURL[key]; ok {
		return fmt.Errorf("key %s is already in database", value)
	}
	// записываем в хранилице userid, id, URL
	ms.IDURL[key] = value
	ms.UserID[key] = userid
	// открываем файл
	sfile, err := os.OpenFile(ms.pathName, os.O_WRONLY, 0777)
	if err != nil {
		log.Println("storage file opening error: ", err)
		return err
	}
	defer sfile.Close()
	// кодирование в JSON
	js, err := json.Marshal(&ms)
	if err != nil {
		log.Println("JSON marshalling from struct error: ", err)
		return err
	}
	// запись в файл
	sfile.Write(js) */
	return nil
}

// конструктор нового хранилища JSON
func NewSQLStorage(u map[string]int, s map[string]string, p string) *StorageSQL {
	db, err := sql.Open("pgx", p)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	return &StorageSQL{
		//UserID:   u,
		//IDURL:    s,
		//pathName: p,

		PostgreSQL: p,
	}
}

// метод получения записи из хранилища
func (ms *StorageSQL) GetFromStorage(key string) (value string, err error) {
	/* value, ok := ms.IDURL[key]
	if !ok {
		return "", fmt.Errorf("key %v not found", key)
	} */
	return value, nil
}

// метод определения длинны хранилища
func (ms *StorageSQL) LenStorage() (lenn int) {

	//lenn = len(ms.IDURL)
	return lenn
}

// метод отбора URLs по UserID
func (ms *StorageSQL) URLsByUserID(userid int) (userURLs map[string]string, err error) {
	/* userURLs = make(map[string]string)
	for k, v := range ms.UserID {
		if v == userid {
			userURLs[k] = ms.IDURL[k]
		}
	}
	if len(userURLs) == 0 {
		err = fmt.Errorf("userid not found in the storage")
	} */
	return userURLs, err
}

func (ms *StorageSQL) LoadFromFileToStorage() {
	// загрузка базы из JSON
	/* p := ms.pathName
	_, pathOk := os.Stat(filepath.Dir(p))
	if os.IsNotExist(pathOk) {
		os.MkdirAll(filepath.Dir(p), 0777)
		log.Printf("folder %s created\n", filepath.Dir(p))
	}
	sfile, err := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("file creating error: ", err)
	}
	defer sfile.Close()

	fileInfo, _ := os.Stat(p)
	if fileInfo.Size() != 0 {
		b, err := io.ReadAll(sfile)
		if err != nil {
			log.Println("file storage reading error:", err)
		}
		err = json.Unmarshal(b, &ms)
		if err != nil {
			log.Println("JSON unmarshalling to struct error:", err)
		}
	} */
}

// посик userid в хранилице
func (ms *StorageSQL) UserIDExist(userid int) bool {
	// цикл по map поиск значения без ключа
	/* for _, v := range ms.UserID {
		if v == userid {
			return true
		}
	} */
	return false
}

func (ms *StorageSQL) StorageOkPing() bool {
 	db, err := sql.Open("pgx", ms.PostgreSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close() 

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return false
	}
	return true
}

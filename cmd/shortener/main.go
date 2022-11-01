package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	//_ "github.com/jackc/pgx/v5/stdlib"
)

// переменные по умолчанию
const (
	defServAddr    = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
	defStoragePath = "db/keyvalue.json"
	defDBlink      = "postgres://postgres:1818@localhost:5432/dbo"
)

// константы цветового вывода в консоль
const (
	colorBlack  = "\u001b[30m"
	colorRed    = "\u001b[31m"
	colorGreen  = "\u001b[32m"
	colorYellow = "\u001b[33m"
	colorBlue   = "\u001b[34m"
	colorReset  = "\u001b[0m"
)

func main() {
	// получаем переменные
	dlink, path, base, addr := flagsVars()
	// инициализируем конструкторы
	s := newStrorageProvider(dlink, path)
	defer s.StorageConnectionClose()
	srvs := services.NewService(s, base)
	h := handlers.NewHandler(srvs, base)
	r := httprouters.NewRouter(h)
	// запускаем сервер
	log.Println("base URL:", colorGreen, base, colorReset)
	log.Println("starting server on:", colorBlue, addr, colorReset)
	log.Fatal(http.ListenAndServe(addr, r))
}

// парсинг флагов и валидация переменных окружения
func flagsVars() (dlink string, path string, base string, addr string) {
	// описываем флаги
	addrFlag := flag.String("a", defServAddr, "HTTP Server address")
	baseFlag := flag.String("b", defBaseURL, "Base URL")
	pathFlag := flag.String("f", defStoragePath, "File storage path")
	dlinkFlag := flag.String("d", "", "Database DSN link")
	// парсим флаги в переменные
	flag.Parse()
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || !govalidator.IsURL(addr) || addr == "" {
		log.Println("eviroment variable SERVER_ADDRESS is empty or has wrong value ", addr)
		addr = *addrFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	base, ok = os.LookupEnv("BASE_URL")
	if !ok || !govalidator.IsURL(base) || base == "" {
		log.Println("eviroment variable BASE_URL is empty or has wrong value ", base)
		base = *baseFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	dlink, ok = os.LookupEnv("DATABASE_DSN")
	if !ok {
		log.Println("eviroment variable DATABASE_DSN is not exist", dlink)
		dlink = *dlinkFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	path, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || (path == "" || !govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ", path)
		path = *pathFlag
	}
	return dlink, path, base, addr
}

// создание интерфейса хранилища
func newStrorageProvider(dlink, path string) (s services.StorageProvider) {
	// если переменная SQL url не пустая, то используем SQL хранилище
	if dlink != "" {
		s = storage.NewSQLStorage(dlink)
		log.Println("server will start with data storage "+colorYellow+"in PostgreSQL:", dlink, colorReset)
		return s
	}
	// иначе если есть path используем для хранения файл
	if path != "" && (govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) {
		log.Println("server will start with data storage " + colorYellow + "in file and memory cash" + colorReset)
		log.Printf("File storage path: %s\n", path)
		s = storage.NewFileStorage(make(map[string]string), make(map[string]string), path)
		s.LoadFromFileToStorage()
		return s
	}
	// если переменная path не валидна, то используем память для хранения id:url
	s = storage.NewMapStorage(make(map[string]string), make(map[string]string))
	log.Println("server will start with data storage " + colorYellow + "in memory" + colorReset)
	return s
}


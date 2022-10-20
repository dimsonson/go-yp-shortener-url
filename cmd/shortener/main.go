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
	_ "github.com/jackc/pgx/v5/stdlib"
)

// переменные по умолчанию
const (
	defServAddr    = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
	defStoragePath = "db/keyvalue.json"
	defDBlink      = "postgres://postgres:1818@localhost:5432/dbo"
)

func main() {
	// описываем флаги
	addrFlag := flag.String("a", defServAddr, "HTTP Server address")
	baseFlag := flag.String("b", defBaseURL, "Base URL")
	pathFlag := flag.String("f", defStoragePath, "File storage path")
	dlinkFlag := flag.String("d", "", "Database DSN link")
	// пасрсим флаги в переменные
	flag.Parse()
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || !govalidator.IsURL(addr) || addr == "" {
		log.Println("eviroment variable SERVER_ADDRESS is empty or has wrong value ", addr)
		addr = *addrFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	base, ok := os.LookupEnv("BASE_URL")
	if !ok || !govalidator.IsURL(base) || base == "" {
		log.Println("eviroment variable BASE_URL is empty or has wrong value ", base)
		base = *baseFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	dlink, ok := os.LookupEnv("DATABASE_DSN")
	if !ok {
		log.Println("eviroment variable DATABASE_DSN is not exist", dlink)
		dlink = *dlinkFlag
	}

	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	path, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || (!govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) || path == "" {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ", path)
		path = *pathFlag
	}

	// задаем переменную провайдера хранилища
	var s services.StorageProvider
	// var d *pgxpool.Pool

	if dlink != "" {
		s = storage.NewSQLStorage(dlink)
		log.Println("server will start with data storage "+ColorYellow+"in PostgreSQL:", dlink, ColorReset)
		defer s.StorageConnectionClose()
		//defer d.Close()
	} else {
		// если переменная не валидна, то используем память для хранения id:url
		if (!govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) || path == "" {
			s = storage.NewMapStorage(make(map[string]int), make(map[string]string))
			log.Println("server will start with data storage" + ColorYellow + "in memory" + ColorReset)
		} else {
			// иначе используем для хранения id:url файл
			s = storage.NewFileStorage(make(map[string]int), make(map[string]string), path)
			s.LoadFromFileToStorage()
			log.Println("server will start with data storage" + ColorYellow + "in file and memory cash" + ColorReset)
			log.Printf("File storage path: %s\n", path)
		}
	}
	// инициализируем конструкторы
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs, base)
	r := httprouters.NewRouter(h)

	// запускаем сервер
	log.Println("base URL:", ColorGreen, base, ColorReset)
	log.Println("starting server on:", ColorBlue, addr, ColorReset)
	log.Fatal(http.ListenAndServe(addr, r))
}

// константы цветового вывода в консоль
const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)

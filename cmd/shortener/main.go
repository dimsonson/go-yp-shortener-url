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
	// парсим флаги в переменные
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
	if !ok || (!govalidator.IsUnixFilePath(path) || !govalidator.IsWinFilePath(path)) || path == "" {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ", path)
		path = *pathFlag
	}

	// задаем переменную провайдера хранилища
	var s services.StorageProvider
	// var d *pgxpool.Pool

	if dlink != "" {
		/* s = storage.NewSQLStorage(dlink)
		log.Println("server will start with data storage "+colorYellow+"in PostgreSQL:", dlink, colorReset)
		defer s.StorageConnectionClose()
		//defer d.Close() */
		//return 
		s = sqlStorageInit(dlink)  
	}

	// если переменная не валидна, то используем память для хранения id:url
	if dlink == "" && (path == "" || (!govalidator.IsUnixFilePath(path) || !govalidator.IsWinFilePath(path))) {
		/* s = storage.NewMapStorage(make(map[string]int), make(map[string]string))
		log.Println("server will start with data storage" + colorYellow + "in memory" + colorReset) */
		s = memoryStrorageInit() 
	}

	if dlink == "" && (path != "" || (govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path))) {
		// иначе используем для хранения id:url файл
		/* s = storage.NewFileStorage(make(map[string]int), make(map[string]string), path)
		s.LoadFromFileToStorage()
		log.Println("server will start with data storage" + colorYellow + "in file and memory cash" + colorReset)
		log.Printf("File storage path: %s\n", path) */
		//return 
		s = fileStrorageInit(path)

	}

	// инициализируем конструкторы
	
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs, base)
	r := httprouters.NewRouter(h)

	// запускаем сервер
	log.Println("base URL:", colorGreen, base, colorReset)
	log.Println("starting server on:", colorBlue, addr, colorReset)
	log.Fatal(http.ListenAndServe(addr, r))
}

// константы цветового вывода в консоль
const (
	colorBlack  = "\u001b[30m"
	colorRed    = "\u001b[31m"
	colorGreen  = "\u001b[32m"
	colorYellow = "\u001b[33m"
	colorBlue   = "\u001b[34m"
	colorReset  = "\u001b[0m"
)

func memoryStrorageInit() services.StorageProvider {
	s := storage.NewMapStorage(make(map[string]int), make(map[string]string))
	log.Println("server will start with data storage" + colorYellow + "in memory" + colorReset)
	return s
}

func sqlStorageInit(dlink string) services.StorageProvider {
	s := storage.NewSQLStorage(dlink)
	log.Println("server will start with data storage "+colorYellow+"in PostgreSQL:", dlink, colorReset)
	defer s.StorageConnectionClose()
	//defer d.Close()
	return s
}

func fileStrorageInit(path string) services.StorageProvider {
	s := storage.NewFileStorage(make(map[string]int), make(map[string]string), path)
	s.LoadFromFileToStorage()
	log.Println("server will start with data storage" + colorYellow + "in file and memory cash" + colorReset)
	log.Printf("File storage path: %s\n", path)
	return s
}

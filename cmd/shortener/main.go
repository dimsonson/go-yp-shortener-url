// Сервис выдачи коротких ссылок по API запросам.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"golang.org/x/crypto/acme/autocert"
)

// Константы по умолчанию.
const (
	defServAddr    = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
	defStoragePath = "db/keyvalue.json"
	defDBlink      = "postgres://postgres:1818@localhost:5432/dbo"
	defHTTPS       = false
)

// Глобальные переменные для использования при сборке - go run -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X main.buildCommit=final"  main.go.
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Вывод данных о версии, дате, коммите сборки.
	log.Printf("version=%s, date=%s, commit=%s\n", buildVersion, buildDate, buildCommit)

	// Получаем переменные из флагов или переменных оркужения в структуру models.Config.
	cfg := flagsVars()
	// Инициализируем конструкторы.
	// Конструктор хранилища.
	s := newStrorageProvider(cfg.DatabaseDsn, cfg.FileStoragePath)
	defer s.Close()
	// Конструктор Put слоя.
	svcRand := &service.Rand{}
	svsPut := service.NewPutService(s, cfg.BaseURL, svcRand)
	hPut := handlers.NewPutHandler(svsPut, cfg.BaseURL)
	// Конструктор Get слоя.
	svsGet := service.NewGetService(s, cfg.BaseURL)
	hGet := handlers.NewGetHandler(svsGet, cfg.BaseURL)
	// Конструктор Delete слоя.
	svsDel := service.NewDeleteService(s, cfg.BaseURL)
	hDel := handlers.NewDeleteHandler(svsDel, cfg.BaseURL)
	// Констуктор Ping слоя.
	svsPing := service.NewPingService(s)
	hPing := handlers.NewPingHandler(svsPing, cfg.BaseURL)
	// Инциализация хендлеров.
	r := httprouters.NewRouter(hPut, hGet, hDel, hPing)

	// Запуск сервера.
	log.Println("base URL:", settings.ColorGreen, cfg.BaseURL, settings.ColorReset)
	if cfg.EnableHTTPS {
		log.Println("starting", settings.ColorBlue, "https", settings.ColorReset, "server on:", settings.ColorBlue, cfg.ServerAddress, settings.ColorReset)
		log.Println(http.Serve(autocert.NewListener(cfg.ServerAddress), r))
		return
	}
	log.Println("starting server on:", settings.ColorBlue, cfg.ServerAddress, settings.ColorReset)
	log.Println(http.ListenAndServe(cfg.ServerAddress, r))
}

// flagsVars парсинг флагов и валидация переменных окружения.
func flagsVars() (cfg models.Config) {
	// описываем флаги
	addrFlag := flag.String("a", "", "HTTP/HTTPS Server address")
	baseFlag := flag.String("b", "", "dase URL")
	pathFlag := flag.String("f", "", "File storage path")
	dlinkFlag := flag.String("d", "", "database DSN link")
	tlsFlag := flag.Bool("s", false, "run HTTPS server")
	cfgFlag := flag.String("c", "", "config json file name")
	// парсим флаги в переменные
	flag.Parse()
	var ok bool
	// используем структуру cfg models.Config для хранения параментров необходимых для запуска сервера
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	cfg.ConfigJSON, ok = os.LookupEnv("CONFIG")
	if !ok && *cfgFlag != "" {
		log.Println("eviroment variable CONFIG is empty or has wrong value ", cfg.ConfigJSON)
		cfg.ConfigJSON = *cfgFlag
	}
	//
	if cfg.ConfigJSON != "" {
		configFile, err := os.ReadFile(*cfgFlag)
		if err != nil {
			log.Println("reading config file error:", err)
		}
		if err == nil {
			err = json.Unmarshal(configFile, &cfg)
			if err != nil {
				log.Printf("unmarshal config file error: %s", err)
			}
		}
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	ServerAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		cfg.ServerAddress = ServerAddress
	}
	if (!ok || !govalidator.IsURL(cfg.ServerAddress) || cfg.ServerAddress == "") && *addrFlag != "" {
		log.Println("eviroment variable SERVER_ADDRESS is empty or has wrong value ")
		cfg.ServerAddress = *addrFlag
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	if *addrFlag == "" {
		cfg.ServerAddress = defServAddr
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	BaseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		cfg.BaseURL = BaseURL
	}
	if (!ok || !govalidator.IsURL(cfg.BaseURL) || cfg.BaseURL == "") && *baseFlag != "" {
		log.Println("eviroment variable BASE_URL is empty or has wrong value ")
		cfg.BaseURL = *baseFlag
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	if *baseFlag == "" {
		cfg.BaseURL = defBaseURL
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	DatabaseDsn, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.DatabaseDsn = DatabaseDsn
	}
	if !ok && *dlinkFlag != "" {
		log.Println("eviroment variable DATABASE_DSN is not exist")
		cfg.DatabaseDsn = *dlinkFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	FileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.FileStoragePath = FileStoragePath
	}
	if !ok || (cfg.FileStoragePath == "" || !govalidator.IsUnixFilePath(cfg.FileStoragePath) || govalidator.IsWinFilePath(cfg.FileStoragePath)) {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ")
		cfg.FileStoragePath = *pathFlag
	}
	// проверяем наличие флага или пременной окружения для старта в https (tls)
	EnableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok || EnableHTTPS == "true" || *tlsFlag {
		cfg.EnableHTTPS = true
		return cfg
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	cfg.EnableHTTPS = defHTTPS
	log.Println("eviroment variable ENABLE_HTTPS is empty or has wrong value ")
	return cfg
}

// newStrorageProvider инциализация интерфейса хранилища в зависимости от переменных окружения и флагов.
func newStrorageProvider(dlink, path string) (s service.StorageProvider) {
	// если переменная SQL url не пустая, то используем SQL хранилище
	if dlink != "" {
		s = storage.NewSQLStorage(dlink)
		log.Println("server will start with data storage "+settings.ColorYellow+"in PostgreSQL:", dlink, settings.ColorReset)
		return s
	}
	// иначе если есть path используем для хранения файл
	if path != "" && (govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) {
		log.Println("server will start with data storage " + settings.ColorYellow + "in file and memory cash" + settings.ColorReset)
		log.Printf("file storage path: %s\n", path)
		s = storage.NewFileStorage(make(map[string]string), make(map[string]string), make(map[string]bool), path)
		s.Load()
		return s
	}
	// если переменная path не валидна, то используем память для хранения id:url
	s = storage.NewMapStorage(make(map[string]string), make(map[string]string), make(map[string]bool))
	log.Println("server will start with data storage " + settings.ColorYellow + "in memory" + settings.ColorReset)
	return s
}

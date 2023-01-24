// Сервис выдачи коротких ссылок по API запросам.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/asaskevich/govalidator"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
)

// Константы по умолчанию.
const (
	defServAddr    = "localhost:443"
	defBaseURL     = "http://localhost:443"
	defStoragePath = "db/keyvalue.json"
	defDBlink      = "postgres://postgres:1818@localhost:5432/dbo"
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

	// Получаем переменные из флагов или переменных оркужения.
	dlink, path, base, addr := flagsVars()

	// Инициализируем конструкторы.
	// Конструктор хранилища.
	s := newStrorageProvider(dlink, path)
	defer s.Close()
	// Конструктор Put слоя.
	svcRand := &service.Rand{}
	svsPut := service.NewPutService(s, base, svcRand)
	hPut := handlers.NewPutHandler(svsPut, base)
	// Конструктор Get слоя.
	svsGet := service.NewGetService(s, base)
	hGet := handlers.NewGetHandler(svsGet, base)
	// Конструктор Delete слоя.
	svsDel := service.NewDeleteService(s, base)
	hDel := handlers.NewDeleteHandler(svsDel, base)
	// Констуктор Ping слоя.
	svsPing := service.NewPingService(s)
	hPing := handlers.NewPingHandler(svsPing, base)
	// Инциализация хендлеров.
	r := httprouters.NewRouter(hPut, hGet, hDel, hPing)

	// Запуск сервера.
	log.Println("base URL:", settings.ColorGreen, base, settings.ColorReset)
	log.Println("starting server on:", settings.ColorBlue, addr, settings.ColorReset)
	//log.Fatal(http.ListenAndServe(addr, r))
	//log.Println(http.Serve(autocert.NewListener(addr), r))
	log.Println(http.ListenAndServeTLS(addr, "cert.pem", "key.pem", r))
	//tls.Listen()
}

// flagsVars парсинг флагов и валидация переменных окружения.
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
		log.Printf("File storage path: %s\n", path)
		s = storage.NewFileStorage(make(map[string]string), make(map[string]string), make(map[string]bool), path)
		s.Load()
		return s
	}
	// если переменная path не валидна, то используем память для хранения id:url
	s = storage.NewMapStorage(make(map[string]string), make(map[string]string), make(map[string]bool))
	log.Println("server will start with data storage " + settings.ColorYellow + "in memory" + settings.ColorReset)
	return s
}

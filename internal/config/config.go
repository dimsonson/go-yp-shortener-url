package config

import (
	"flag"
	"log"
	"os"

	"github.com/asaskevich/govalidator"
)

// переменные по умолчанию
var (
	ServAddr    = "localhost:8080"
	BaseURL     = "http://localhost:8080"
	StoragePath = "db/keyvalue.json"
)

func config() {
	// описываем флаги
	flag.StringVar(&ServAddr, "a", ServAddr, "HTTP Server address")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base URL")
	flag.StringVar(&StoragePath, "f", StoragePath, "File storage path")
	// пасрсим флаги в переменные
	flag.Parse()
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	Addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok || govalidator.IsURL(Addr) || Addr != "" {
		log.Println("eviroment variable SERVER_ADDRESS is empty or has wrong value ", Addr)
		ServAddr = Addr //= *addrFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	Base, ok := os.LookupEnv("BASE_URL")
	if ok || govalidator.IsURL(Base) || Base != "" {
		log.Println("eviroment variable BASE_URL is empty or has wrong value ", Base)
		BaseURL = Base //= *baseFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	Path, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok || (govalidator.IsUnixFilePath(Path) || govalidator.IsWinFilePath(Path)) || Path != "" {
		log.Println("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ", Path)
		StoragePath = Path //= *pathFlag

	}
}

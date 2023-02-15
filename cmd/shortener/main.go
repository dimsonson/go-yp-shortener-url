// Сервис выдачи коротких ссылок по API запросам.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006/01/02 15:04:05"})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	//log.Logger := zerolog.New(os.Stderr)
}

// Глобальные переменные для использования при сборке - go run -ldflags "-X main.buildVersion=v0.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X main.buildCommit=final"  main.go.
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	// Вывод данных о версии, дате, коммите сборки.
	log.Printf("version=%s, date=%s, commit=%s", buildVersion, buildDate, buildCommit)

	// опередяляем контекст уведомления о сигнале прерывания
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// создание конфигурацию сервера
	cfg := server.NewConfig()
	// парсинг конфигурации сервера
	cfg.Parse()
	// создание сервера
	srv := server.NewServer(ctx, stop, *cfg)
	// запуск сервера
	srv.Start()
	// остановка всех сущностей, куда передан контекст по прерыванию
	stop()
	// ожидаем выполнение горутин
	srv.Wg.Wait()
	// логирование закрытия сервера без ошибок
	log.Print("http/grpc server gracefully shutdown")
}

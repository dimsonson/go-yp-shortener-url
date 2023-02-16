// server пакет конфигурирования, запуска, остановки серверов HTTP или GRPC в заисимости от конфигурации.
package server

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/asaskevich/govalidator"
	pb "github.com/dimsonson/go-yp-shortener-url/internal/app/api/grpc/proto"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Константы по умолчанию.
const (
	defServAddr    = "localhost:8080"
	defBaseURL     = "http://localhost:8080"
	defStoragePath = "db/keyvalue.json"
	defDBlink      = "postgres://postgres:1818@localhost:5432/dbo"
	defHTTPS       = false
)

// Server структура для хранения серверов.
type Server struct {
	GRPCserver *grpc.Server
	HTTPserver *http.Server
	Config
	Wg   sync.WaitGroup
	Ctx  context.Context
	Stop context.CancelFunc
	//PutServer
	PutServicePrivider
	PutService *PutServices
}

// Config структура конфигурации сервиса, при запуске сервиса с флагом -c/config
// и отсутствии иных флагов и переменных окружения заполняется из файла указанного в этом флаге или переменной окружения CONFIG.
type Config struct {
	ServerAddress   string     `json:"server_address"`
	BaseURL         string     `json:"base_url"`
	FileStoragePath string     `json:"file_storage_path"`
	DatabaseDsn     string     `json:"database_dsn"`
	EnableHTTPS     bool       `json:"enable_https"`
	TrustedSubnet   string     `json:"trusted_subnet"`
	EnableGRPC      bool       `json:"enable_grpc"`
	TrustedCIDR     *net.IPNet `json:"-"`
	ConfigJSON      string     `json:"-"`
}

// NewServer конструктор создания нового сервера в соответствии с существующей конфигурацией.
func NewServer(ctx context.Context, stop context.CancelFunc, cfg Config) *Server {
	return &Server{
		Config: cfg,
		Ctx:    ctx,
		Stop:   stop,
	}
}

// NewConfig конструктор создания конфигурации сервера из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func NewConfig() *Config {
	return &Config{}
}

// Parse метод парсинга и получения значений из переменных оружения, флагов, конфиг файла, а так же значений по умолчанию.
func (cfg *Config) Parse() {
	// описываем флаги
	addrFlag := flag.String("a", "", "HTTP/HTTPS Server address")
	baseFlag := flag.String("b", "", "dase URL")
	pathFlag := flag.String("f", "", "File storage path")
	dlinkFlag := flag.String("d", "", "database DSN link")
	tlsFlag := flag.Bool("s", false, "run as HTTPS server")
	cfgFlag := flag.String("c", "", "config json file name")
	trustFlag := flag.String("t", "", "trusted subnet CIDR for /api/internal/stats")
	// парсим флаги в переменные
	flag.Parse()
	var ok bool
	// используем структуру cfg models.Config для хранения параментров необходимых для запуска сервера
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	cfg.ConfigJSON, ok = os.LookupEnv("CONFIG")
	if !ok && *cfgFlag != "" {
		log.Print("eviroment variable CONFIG is empty or has wrong value ", cfg.ConfigJSON)
		cfg.ConfigJSON = *cfgFlag
	}
	// читаем конфигурвационный файл и парксим в стркутуру
	if cfg.ConfigJSON != "" {
		configFile, err := os.ReadFile(*cfgFlag)
		if err != nil {
			log.Print("reading config file error:", err)
		}
		if err == nil {
			err = json.Unmarshal(configFile, &cfg)
			if err != nil {
				log.Printf("unmarshal config file error: %s", err)
			}
		}
	}
	// проверяем наличие флага или пременной окружения для CIDR доверенной сети эндпойнта /api/internal/stats
	TrustedSubnet, ok := os.LookupEnv("TRUSTED_SUBNET")
	if ok {
		cfg.TrustedSubnet = TrustedSubnet
	}
	if *trustFlag != "" {
		cfg.TrustedSubnet = *trustFlag
	}
	if cfg.TrustedSubnet != "" {
		var err error
		_, cfg.TrustedCIDR, err = net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			log.Print("parse CIDR error: ", err)
		}
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	ServerAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if ok {
		cfg.ServerAddress = ServerAddress
	}
	if (!ok || !govalidator.IsURL(cfg.ServerAddress) || cfg.ServerAddress == "") && *addrFlag != "" {
		log.Print("eviroment variable SERVER_ADDRESS is empty or has wrong value ")
		cfg.ServerAddress = *addrFlag
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	if !ok && *addrFlag == "" {
		cfg.ServerAddress = defServAddr
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	BaseURL, ok := os.LookupEnv("BASE_URL")
	if ok {
		cfg.BaseURL = BaseURL
	}
	if (!ok || !govalidator.IsURL(cfg.BaseURL) || cfg.BaseURL == "") && *baseFlag != "" {
		log.Print("eviroment variable BASE_URL is empty or has wrong value ")
		cfg.BaseURL = *baseFlag
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	if !ok && *baseFlag == "" {
		cfg.BaseURL = defBaseURL
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	DatabaseDsn, ok := os.LookupEnv("DATABASE_DSN")
	if ok {
		cfg.DatabaseDsn = DatabaseDsn
	}
	if !ok && *dlinkFlag != "" {
		log.Print("eviroment variable DATABASE_DSN is not exist")
		cfg.DatabaseDsn = *dlinkFlag
	}
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	FileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		cfg.FileStoragePath = FileStoragePath
	}
	if !ok || (cfg.FileStoragePath == "" || !govalidator.IsUnixFilePath(cfg.FileStoragePath) || govalidator.IsWinFilePath(cfg.FileStoragePath)) {
		log.Print("eviroment variable FILE_STORAGE_PATH is empty or has wrong value ")
		cfg.FileStoragePath = *pathFlag
	}
	// проверяем наличие флага или пременной окружения для старта в https (tls)
	EnableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok || EnableHTTPS == "true" || *tlsFlag {
		cfg.EnableHTTPS = true
		return
	}
	// если нет флага или переменной окружения используем переменную по умолчанию
	cfg.EnableHTTPS = defHTTPS
	log.Print("eviroment variable ENABLE_HTTPS is empty or has wrong value ")
}

// Start метод запуска сервара, вид запвсукаемого сервера зависит от EnableGRPC в структуре Config.
func (srv *Server) Start() {
	if srv.EnableGRPC {
		srv.InitGRPC()
		srv.InitGRPCservice()
		//	srv.InitGRPC()
		srv.grpcGracefullShotdown()
		srv.StartGRPC()
		return
	}
	srv.InitHTTP()
	srv.httpGracefullShotdown()
	srv.StartHTTP()
}

// StartHTTPS запуск HTTP сервера.
func (srv *Server) StartHTTP() {
	// Запуск сервера.
	log.Print("base URL:", settings.ColorGreen, srv.BaseURL, settings.ColorReset)
	// Выбор варианта запуска сервера http или https.
	if srv.EnableHTTPS {
		log.Print("starting", settings.ColorBlue, "https", settings.ColorReset, "server on:", settings.ColorBlue, srv.ServerAddress, settings.ColorReset)
		log.Print(srv.HTTPserver.Serve(autocert.NewListener()))
	} else {
		log.Print("starting server on:", settings.ColorBlue, srv.ServerAddress, settings.ColorReset)
		log.Print(srv.HTTPserver.ListenAndServe())
	}
}

// InitHTTP инциализация HTTP сервера.
func (srv *Server) InitHTTP() {
	// Инициализируем конструкторы.
	// Конструктор хранилища.
	s := newStrorageProvider(srv.DatabaseDsn, srv.FileStoragePath)

	// Конструкторы.
	//svcRand := &service.Rand{}
	svsPut := service.NewPutService(s, srv.BaseURL)
	hPut := handlers.NewPutHandler(svsPut, srv.BaseURL)
	// Конструктор Get слоя.
	svsGet := service.NewGetService(s, srv.BaseURL)
	hGet := handlers.NewGetHandler(svsGet, srv.BaseURL)
	// Конструктор Delete слоя.
	svsDel := service.NewDeleteService(s, srv.BaseURL)
	hDel := handlers.NewDeleteHandler(svsDel, srv.BaseURL)
	// Констуктор Ping слоя.
	svsPing := service.NewPingService(s)
	hPing := handlers.NewPingHandler(svsPing, srv.TrustedCIDR)
	// Инциализация хендлеров.
	r := httprouters.NewRouter(hPut, hGet, hDel, hPing)

	// конфигурирование http сервера
	srv.HTTPserver = &http.Server{Addr: srv.ServerAddress, Handler: r}
}

type PutServices struct {
	svsPut  *service.PutServices
	svsGet  *service.GetServices
	svsDel  *service.DeleteServices
	svsPing *service.PingServices
	pb.UnimplementedPutServer
	Server
}

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPC() {
	srv.PutService = &PutServices{}
	// Инициализируем конструкторы.
	// Конструктор хранилища.
	//s := newStrorageProvider(srv.DatabaseDsn, srv.FileStoragePath)

	//fmt.Println(srv.PutServ(srv.Ctx))

	// Конструкторы.
	//svcRand := &service.Rand{}
	//srv.svsPut = service.NewPutService(s, srv.BaseURL) //, svcRand)

	//fmt.Println(srv.svsPut.Put(srv.Ctx, "888", "999"))

	// Конструктор Get слоя.
	//srv.svsGet = service.NewGetService(s, srv.BaseURL)
	// Конструктор Delete слоя.
	//srv.svsDel = service.NewDeleteService(s, srv.BaseURL)
	// Констуктор Ping слоя.
	//srv.svsPing = service.NewPingService(s)

	// создаём gRPC-сервер без зарегистрированной службы
	// srv.GRPCserver = grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(
	// 		logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(log.Logger)),
	// 	),
	// 	//grpc_recovery.UnaryServerInterceptor(),
	// )

}

// InitGRPC инциализация GRPC сервера.
func (svs *Server) InitGRPCservice() {
	// Инициализируем конструкторы.
	// Конструктор хранилища.
	s := newStrorageProvider(svs.DatabaseDsn, svs.FileStoragePath)

	//fmt.Println(s.Len(srv.Ctx))

	// Конструкторы.
	//svcRand := &service.Rand{}
	svs.PutService.svsPut = service.NewPutService(s, svs.BaseURL) //, svcRand)

	//fmt.Println(srv.svsPut.Put(srv.Ctx, "888", "999"))

	// Конструктор Get слоя.
	svs.PutService.svsGet = service.NewGetService(s, svs.BaseURL)
	// Конструктор Delete слоя.
	svs.PutService.svsDel = service.NewDeleteService(s, svs.BaseURL)
	// Констуктор Ping слоя.
	svs.PutService.svsPing = service.NewPingService(s)

	// Define customfunc to handle panic
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	// создаём gRPC-сервер без зарегистрированной службы
	svs.GRPCserver = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(log.Logger)),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		
	)

}

type PutServicePrivider interface {
	InitGRPCservice()
}

// type PutServer struct{
// 	pb.UnimplementedPutServer
// }

// StartGRPC запуск GRPC сервера.
func (srv *Server) StartGRPC() {

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {

	}

	//PutServer := &PutServices{}
	pb.RegisterPutServer(srv.GRPCserver, srv.PutService)

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := srv.GRPCserver.Serve(listen); err != nil {

	}

}

func (s *PutServices) Put(ctx context.Context, rq *pb.PutRequest) (*pb.PutResponse, error) {

	//	st := newStrorageProvider(s.DatabaseDsn, s.FileStoragePath)

	//fmt.Println(s.Len(srv.Ctx))

	// Конструкторы.
	//svcRand := &service.Rand{}
	//	svsPut := service.NewPutService(st, s.BaseURL) //, svcRand)
	var response pb.PutResponse
	//ctx = context.Background()
	key, _ := s.svsPut.Put(ctx, rq.Value, rq.Userid)
	fmt.Println(key)
	response.ExistKey = key //"358ksdHJ" //key //s.svsPut.Put()
	return &response, nil   //err
}

func (s *PutServices) PutBatch(ctx context.Context, rq *pb.PutBatchRequest) (*pb.PutBatchResponse, error) {
	var response pb.PutBatchResponse
	response.DcCorr.CorrelationID = "3GSFHJY"
	//s.svsPut.Put()
	return &response, nil
}

// grpcGracefullShotdown метод благопроиятного для соединений и незавершенных запросов закрытия сервера.
func (srv *Server) grpcGracefullShotdown() {
	srv.Wg.Add(1)
	go func() {
		// получаем сигнал о завершении приложения
		<-srv.Ctx.Done()
		log.Printf("got signal, attempting graceful shutdown")
		//srv.Stop()
		srv.GRPCserver.GracefulStop()
		// grpc.Stop() // leads to error while receiving stream response: rpc error: code = Unavailable desc = transport is closing
		srv.Wg.Done()
	}()
}

// GracefullShotdown метод благопроиятного для соединений и незавершенных запросов закрытия сервера.
func (srv *Server) httpGracefullShotdown() {
	// добавляем счетчик горутины
	srv.Wg.Add(1)
	// запуск горутины shutdown http сервера
	go httpServerShutdown(srv.Ctx, &srv.Wg, srv.HTTPserver)
}

// newStrorageProvider инциализация интерфейса хранилища в зависимости от переменных окружения и флагов.
func newStrorageProvider(dlink, path string) (s service.StorageProvider) {
	// если переменная SQL url не пустая, то используем SQL хранилище
	if dlink != "" {
		s = storage.NewSQLStorage(dlink)
		log.Print("server will start with data storage "+settings.ColorYellow+"in PostgreSQL:", dlink, settings.ColorReset)
		return s
	}
	// иначе если есть path используем для хранения файл
	if path != "" && (govalidator.IsUnixFilePath(path) || govalidator.IsWinFilePath(path)) {
		log.Print("server will start with data storage " + settings.ColorYellow + "in file and memory cash" + settings.ColorReset)
		log.Printf("file storage path: %s\n", path)
		s = storage.NewFileStorage(make(map[string]string), make(map[string]string), make(map[string]bool), path)
		s.Load()
		return s
	}
	// если переменная path не валидна, то используем память для хранения id:url
	s = storage.NewMapStorage(make(map[string]string), make(map[string]string), make(map[string]bool))
	log.Print("server will start with data storage " + settings.ColorYellow + "in memory" + settings.ColorReset)
	return s
}

// httpServerShutdown реализует gracefull shutdown для ListenAndServe.
func httpServerShutdown(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	// получаем сигнал о завершении приложения
	<-ctx.Done()
	// завершаем открытые соединения и закрываем http server
	if err := srv.Shutdown(ctx); err != nil {
		// логирование ошибки остановки сервера
		log.Printf("HTTP server Shutdown error: %v", err)
	}
	// уменьшаем счетчик запущенных горутин
	wg.Done()
}

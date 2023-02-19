// server пакет конфигурирования, запуска, остановки серверов HTTP или GRPC в заисимости от конфигурации.
package server

import (
	"context"
	"encoding/json"
	"flag"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/asaskevich/govalidator"
	pb "github.com/dimsonson/go-yp-shortener-url/internal/app/api/grpc/proto"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/models"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	grpczerolog "github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/jackc/pgerrcode"

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
	ShortServicePrivider
	ShortService *ShortServices
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
	grpcFlag := flag.Bool("g", false, "run as GRPC server")
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
	// проверяем наличие переменной окружения, если ее нет или она не валидна, то используем значение из флага
	EnableGRPC, ok := os.LookupEnv("ENABLE_GRPC")
	if ok && EnableGRPC == "true" || *grpcFlag {
		cfg.EnableGRPC = true
	}
	if !ok {
		log.Print("eviroment variable ENABLE_GRPC is empty or has wrong value ")
	}
	// проверяем наличие флага или пременной окружения для старта в https (tls)
	EnableHTTPS, ok := os.LookupEnv("ENABLE_HTTPS")
	if ok && EnableHTTPS == "true" || *tlsFlag {
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
		//srv.InitGRPC()
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
	svsPut := service.NewPutService(s, srv.BaseURL) //, svcRand)
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

type ShortServices struct {
	svsPut  *service.PutServices
	svsGet  *service.GetServices
	svsDel  *service.DeleteServices
	svsPing *service.PingServices
	pb.UnimplementedShortServiceServer
	Server
}

// InitGRPC инциализация GRPC сервера.
func (srv *Server) InitGRPCservice() {
	// Инициализируем конструкторы.
	srv.ShortService = &ShortServices{}
	// Конструктор хранилища.
	s := newStrorageProvider(srv.DatabaseDsn, srv.FileStoragePath)

	//fmt.Println(s.Len(srv.Ctx))

	// Конструкторы.
	// svcRand := &service.Rand{}
	srv.ShortService.svsPut = service.NewPutService(s, srv.BaseURL) //, svcRand)

	//fmt.Println(srv.srvPut.Put(srv.Ctx, "888", "999"))

	// Конструктор Get слоя.
	srv.ShortService.svsGet = service.NewGetService(s, srv.BaseURL)
	// Конструктор Delete слоя.
	srv.ShortService.svsDel = service.NewDeleteService(s, srv.BaseURL)
	// Констуктор Ping слоя.
	srv.ShortService.svsPing = service.NewPingService(s)

	// Define customfunc to handle panic
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []recovery.Option{
		recovery.WithRecoveryHandler(customFunc),
	}

	// создаём gRPC-сервер без зарегистрированной службы
	srv.GRPCserver = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(grpczerolog.InterceptorLogger(log.Logger)),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)

}

// ShortServicePrivider интерфейс реализации паттерна компоновщик для использования сервисов в серверах.
type ShortServicePrivider interface {
	InitGRPCservice()
}

// StartGRPC запуск GRPC сервера.
func (srv *Server) StartGRPC() {

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("gRPC listener error: %v", err)
	}

	//PutServer := &PutServices{}
	pb.RegisterShortServiceServer(srv.GRPCserver, srv.ShortService)

	log.Print("Сервер gRPC начинает работу")
	// получаем запрос gRPC
	if err := srv.GRPCserver.Serve(listen); err != nil {
		log.Printf("gRPC server error: %v", err)
	}
}

func (s *ShortServices) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	var out pb.PutResponse
	var err error
	out.Key, err = s.svsPut.Put(ctx, in.Value, in.Userid)
	if err != nil {
		log.Printf("call Put error: %v", err)
	}
	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		err = status.Errorf(codes.AlreadyExists, `this link already shortened: %s`, in.Value)
		out.Error = codes.AlreadyExists.String()
	case err != nil:
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	default:
		out.Error = codes.OK.String()
	}
	return &out, err
}

func (s *ShortServices) PutBatch(ctx context.Context, in *pb.PutBatchRequest) (*pb.PutBatchResponse, error) {
	var out pb.PutBatchResponse
	var dcc []models.BatchRequest

	tmpIn := models.BatchRequest{}
	for _, v := range in.Dcc {
		tmpIn.CorrelationID = v.CorrelationID
		tmpIn.OriginalURL = v.OriginalURL
		tmpIn.ShortURL = v.ShortURL
		dcc = append(dcc, tmpIn)
	}

	dcCorr, err := s.svsPut.PutBatch(ctx, dcc, in.Userid)
	tmpOut := pb.BatchResponse{}
	for _, v := range dcCorr {
		tmpOut.CorrelationID = v.CorrelationID
		tmpOut.ShortURL = v.ShortURL
		out.DcCorr = append(out.DcCorr, &tmpOut)
	}

	switch {
	case err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation):
		err = status.Error(codes.AlreadyExists, `one of this this links already shortened`)
		out.Error = codes.AlreadyExists.String()
	case err != nil:
		log.Printf("call Put error: %v", err)
		status.Errorf(codes.Internal, `server error %s`, error.Error(err))
		out.Error = codes.Internal.String()
	default:
		out.Error = codes.OK.String()
	}
	return &out, err
}

// grpcGracefullShotdown метод благопроиятного для соединений и незавершенных запросов закрытия сервера.
func (srv *Server) grpcGracefullShotdown() {
	srv.Wg.Add(1)
	go func() {
		// получаем сигнал о завершении приложения
		<-srv.Ctx.Done()
		log.Print("got signal, attempting graceful shutdown")
		srv.GRPCserver.GracefulStop()
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

package internal

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/shettyh/threadpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Envs        = InitEnvVars()
	SysStatus   = InitSystemStatus()
	Log         = TempLog{}
	Logger      = NewLogger()
	Db          = NewSQLDB()
	redisClient = NewConnectRedis()
)

var (
	Domains           = InitAPIDomains(Envs.IsProduction)
	Endpoints         = InitAPIEndpoints()
	Keys              = InitKeys(Envs.IsProduction)
	Pool              = NewThreadPool()
	Brokers           = NewBrokers(Envs.IsProduction)
	KafkaTopicNameAll = NewKafkaTopicNameAll(Envs.IsProduction)
	ServiceName       = "mine"
	RequestMethod     = httpRequestMethod
)

type HTTPRequestMethod struct {
	GET   string
	POST  string
	PATCH string
	PUT   string
}

var httpRequestMethod = HTTPRequestMethod{
	GET:   "GET",
	POST:  "POST",
	PATCH: "PATCH",
	PUT:   "PUT",
}

const (
	INF_COST                = 1e9
	NEG_INF_COST            = -1e9
	CODE_DB_FAILED          = 300
	CODE_WRONG_PARAMS       = 400
	CODE_RATE_LIMIT         = 401
	CODE_TOKEN_REQUIRED     = 1003
	CODE_TOKEN_EXPIRED      = 1001
	CODE_INVALID_TOKEN      = 1002
	CODE_TOKEN_APP_REQUIRED = 2003
	CODE_TOKEN_APP_EXPIRED  = 2001
	CODE_INVALID_TOKEN_APP  = 2002
	CODE_SYSTEM_BUSY        = 300
	CODE_SYSTEM_ERROR       = 301

	MSG_DB_FAILED          = "Kết nối hệ thống lỗi, vui lòng thử lại sau ít phút"
	MSG_WRONG_PARAMS       = "Sai thông tin đầu vào, vui lòng kiểm tra lại thông tin"
	MSG_RATE_LIMIT         = "Số lần truy cập đã đến giới hạn, vui lòng thử lại sau ít phút."
	MSG_TOKEN_REQUIRED     = "Token không tồn tại"
	MSG_TOKEN_EXPIRED      = "Token hết hạn"
	MSG_INVALID_TOKEN      = "Token không hợp lệ"
	MSG_TOKEN_APP_REQUIRED = "Token không tồn tại"
	MSG_TOKEN_APP_EXPIRED  = "Token hết hạn"
	MSG_INVALID_TOKEN_APP  = "Token không hợp lệ"
	MSG_SYSTEM_BUSY        = "Hệ thống đang bận, vui lòng thử lại sau ít phút"
	MSG_SYSTEM_ERROR       = "Có lỗi trong quá trình xử lý, vui lòng thử lại sau ít phút"
)

type AppKeys struct {
	SampleKey string
}

type EnvVars struct {
	SqlHost      string `mapstructure:"DB_HOST"`
	SqlUser      string `mapstructure:"DB_USERNAME"`
	SqlPassword  string `mapstructure:"DB_PASSWORD"`
	SqlDBName    string `mapstructure:"DB_DATABASE"`
	SqlPort      int    `mapstructure:"DB_PORT"`
	IsProduction bool   `mapstructure:"USE_PRODUCTION"`
	IsDev        bool   `mapstructure:"IS_DEV"`
}
type ApiDomains struct {
	TypeSense string
}

func InitSystemStatus() *AllSystemStatus {
	return &AllSystemStatus{
		DbFailed: &SystemStatus{
			Status: CODE_DB_FAILED,
			Msg:    MSG_DB_FAILED,
		},
		WrongParams: &SystemStatus{
			Status: CODE_WRONG_PARAMS,
			Msg:    MSG_WRONG_PARAMS,
		}, RateLimit: &SystemStatus{
			Status: CODE_RATE_LIMIT,
			Msg:    MSG_RATE_LIMIT,
		}, TokenRequired: &SystemStatus{
			Status: CODE_TOKEN_REQUIRED,
			Msg:    MSG_TOKEN_REQUIRED,
		}, TokenExpired: &SystemStatus{
			Status: CODE_TOKEN_EXPIRED,
			Msg:    MSG_TOKEN_EXPIRED,
		}, InvalidToken: &SystemStatus{
			Status: CODE_INVALID_TOKEN,
			Msg:    MSG_INVALID_TOKEN,
		}, SystemBusy: &SystemStatus{
			Status: CODE_SYSTEM_BUSY,
			Msg:    MSG_SYSTEM_BUSY,
		}, SystemError: &SystemStatus{
			Status: CODE_SYSTEM_ERROR,
			Msg:    MSG_SYSTEM_ERROR,
		}, TokenAppRequired: &SystemStatus{
			Status: CODE_TOKEN_APP_REQUIRED,
			Msg:    MSG_TOKEN_APP_REQUIRED,
		}, InvalidTokenApp: &SystemStatus{
			Status: CODE_INVALID_TOKEN_APP,
			Msg:    MSG_INVALID_TOKEN_APP,
		}, TokenAppExpired: &SystemStatus{
			Status: CODE_TOKEN_APP_EXPIRED,
			Msg:    MSG_TOKEN_APP_EXPIRED,
		},
	}
}

type SystemStatus struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Detail interface{} `json:"detail"`
}
type AllSystemStatus struct {
	DbFailed         *SystemStatus
	WrongParams      *SystemStatus
	RateLimit        *SystemStatus
	TokenRequired    *SystemStatus
	TokenExpired     *SystemStatus
	InvalidToken     *SystemStatus
	SystemBusy       *SystemStatus
	SystemError      *SystemStatus
	TokenAppRequired *SystemStatus
	InvalidTokenApp  *SystemStatus
	TokenAppExpired  *SystemStatus
}

type TypeSenseEndpoint struct {
	TextSearch string
}
type ApiEndpoints struct {
	TypeSense TypeSenseEndpoint
}

func NewThreadPool() *threadpool.ThreadPool {
	fmt.Println("LOADING THREAD POOL ...")
	threadPool := threadpool.NewThreadPool(50, 100000)
	if threadPool == nil {
		fmt.Println("Failed")
	}
	fmt.Println("LOADING THREAD POOL SUCCESS...")
	return threadPool
}

func InitEnvVars() *EnvVars {
	fmt.Println("LOADING ENVS...")
	envs := &EnvVars{}
	viper.SetConfigFile("app.env")
	errEnvFile := viper.ReadInConfig()
	if errEnvFile != nil {
		viper.AutomaticEnv()
		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_USERNAME")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_DATABASE")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("USE_PRODUCTION")
		viper.BindEnv("IS_DEV")
	}
	if err := viper.Unmarshal(envs); err != nil {
		fmt.Println("Error viper.Unmarshal", err)
		fmt.Println("LOADING ENVS FAILED")
	}
	fmt.Println("LOADING ENVS SUCCESS")
	return envs
}

func NewSQLDB() *gorm.DB {
	fmt.Println("LOADING MYSQL DB 1...")
	DBDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
		Envs.SqlUser, Envs.SqlPassword, Envs.SqlHost, Envs.SqlPort, Envs.SqlDBName)
	DBDSN = DBDSN + "&loc=Asia%2FHo_Chi_Minh"
	DB, err := sql.Open("mysql", DBDSN)
	if err != nil {
		fmt.Println("NewSQLDB", err)
		return nil
	}
	DB.SetConnMaxLifetime(time.Minute * 10)
	DB.SetMaxOpenConns(10000)
	DB.SetMaxIdleConns(1000)
	errPing := DB.Ping()
	if errPing != nil {
		fmt.Println("DB PING ERR: ", errPing)
		return nil
	}
	customLogger := &CustomLogger{
		Config: logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DB,
	}), &gorm.Config{
		Logger: customLogger,
	})
	if err != nil {
		fmt.Println("gorm DB connect failed", err.Error())
		return nil
	}
	fmt.Println("LOADING MYSQL DB SUCCESS...")
	gormDB.Debug()
	return gormDB
}

func CheckSQLDB() (*gorm.DB, error) {
	fmt.Println("LOADING MYSQL DB ...")
	DBDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
		Envs.SqlUser, Envs.SqlPassword, Envs.SqlHost, Envs.SqlPort, Envs.SqlDBName)
	DBDSN = DBDSN + "&loc=Asia%2FHo_Chi_Minh"
	DB, err := sql.Open("mysql", DBDSN)
	if err != nil {
		// cm.Log.Debug("Open mysql connection failed", zap.Error(err))
		fmt.Println("NewSQLDB", err)
		return nil, err
	}
	DB.SetConnMaxLifetime(time.Minute * 10)
	DB.SetMaxOpenConns(10000)
	DB.SetMaxIdleConns(1000)
	errPing := DB.Ping()
	if errPing != nil {
		fmt.Println("DB PING ERR: ", errPing)
		return nil, err
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DB,
	}))
	if err != nil {
		fmt.Println("gorm DB connect failed", err.Error())
		return nil, err
	}
	fmt.Println("LOADING MYSQL DB SUCCESS...")
	return gormDB, nil
}

func NewConnectRedis() *redis.Client {
	fmt.Println("NewConnectRedis start...")
	client := redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	fmt.Println("NewConnectRedis SUCCESS")
	return client
}

func NewBrokers(useProduction bool) []string {
	if Envs.IsDev {
		return []string{"kafka-1:19092", "kafka-2:29092", "kafka-3:39092"}
	}
	if useProduction {
		return []string{"isc-kafka01:9092", "isc-kafka02:9092", "isc-kafka03:9092"}
	}
	return []string{"isc-kafka01:9092", "isc-kafka02:9092", "isc-kafka03:9092"}
}

func NewKafkaTopicNameAll(useProduction bool) string {
	if Envs.IsDev {
		return "dev-mine"
	}
	if useProduction {
		return "mine"
	}
	return "stag-mine"
}

func InitAPIDomains(isProduction bool) *ApiDomains {
	// production
	if isProduction {
		return &ApiDomains{}
	}
	// staging
	return &ApiDomains{
		TypeSense: "http://typesense-stag/",
	}
}

func InitAPIEndpoints() *ApiEndpoints {
	endpoints := &ApiEndpoints{
		TypeSense: TypeSenseEndpoint{
			TextSearch: "/collections/{{document}}/documents/search",
		},
	}
	return endpoints
}

func InitKeys(isProduction bool) *AppKeys {
	keys := &AppKeys{
		SampleKey: "sample_key",
	}
	if isProduction {
	}
	return keys
}

func NewLogger() *zap.Logger {
	// Log vào stdOut
	writeSyncer := zapcore.AddSync(os.Stdout)
	// Thiết lập log
	loggerCore := zapcore.NewCore(logEndcoder(), writeSyncer, zap.DebugLevel)
	return zap.New(loggerCore)
}

func logEndcoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = zapcore.TimeEncoderOfLayout("06-01-02 15:04:05")
	return zapcore.NewJSONEncoder(encodeConfig)
}

type TempLog struct{}

func (TempLog) Info(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	fields = append(fields, zap.Any("caller", fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)))
	Logger.Info(msg, fields...)
}

func (TempLog) InfoQuery1(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(4) // dòng repo/query
	fn := runtime.FuncForPC(pc)
	pc1, file1, line1, _ := runtime.Caller(5) // dòng call hàm repo nếu ko tạo transaction
	pc2, file2, line2, _ := runtime.Caller(7) // Nếu có transaction thì là dòng call hàm repo
	listString := strings.Split(file, ServiceName+"/internal/repositories/")
	fields = append(fields, zap.Any("caller", fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), listString[len(listString)-1], line)))
	caller1 := ""
	func1 := ""
	callerLine1 := 0
	if strings.Contains(file1, ServiceName) {
		func1 = shortFuncNameCaller(runtime.FuncForPC(pc1).Name())
		caller1 = shortCaller(file1)
		callerLine1 = line1
	} else {
		func1 = runtime.FuncForPC(pc2).Name()
		caller1 = shortCaller(file2)
		callerLine1 = line2
	}
	fields = append(fields, zap.Any("caller1", fmt.Sprintf("%v %v:%v", func1, caller1, callerLine1)))
	Logger.Info(msg, fields...)
}

func (TempLog) InfoQuery(msg string, query zap.Field, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(4)
	fn := runtime.FuncForPC(pc)
	callerStr := fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), shortCaller(file), line)
	if strings.Contains(file, "preload") {
		fmt.Printf("%v--PRELOAD---%v\n", msg, query.String)
	} else if strings.Contains(file, "associations") {
		fmt.Printf("%v--ASSOCIATION---%v\n", msg, query.String)
	} else {
		pc1, file1, line1, _ := runtime.Caller(5) // dòng call hàm repo nếu ko tạo transaction
		caller1Str := ""
		if !strings.Contains(file1, ServiceName) {
			pc1, file1, line1, _ = runtime.Caller(7) // Nếu có transaction thì là dòng call hàm repo
		}
		caller1Str = fmt.Sprintf("%v %v:%v", shortFuncNameCaller(runtime.FuncForPC(pc1).Name()), shortCaller(file1), line1)
		fmt.Printf("%v---%v\n%v----%v\n\n", msg, query.String, callerStr, caller1Str)
		fields = append(fields, zap.Any("caller1", caller1Str))
		fields = append(fields, zap.Any("caller", callerStr))
	}
}

func (TempLog) Error(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	fields = append(fields, zap.Any("caller", fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)))
	Logger.Error(msg, fields...)
}

func (TempLog) Debug(msg string, fields ...zap.Field) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	listString := strings.Split(file, ServiceName+"/")
	file = shortCaller(listString[len(listString)-1])
	fields = append(fields, zap.Any("caller", fmt.Sprintf("%v %v:%v", shortFuncNameCaller(fn.Name()), file, line)))
	Logger.Debug(msg, fields...)
}

func shortCaller(caller string) string {
	if strings.Contains(caller, "cmd/") {
		caller = strings.Split(caller, "cmd/")[1]
	}
	if strings.Contains(caller, "internal/") {
		caller = strings.Split(caller, "internal/")[1]
	}
	if strings.Contains(caller, "delivery/") {
		caller = strings.Split(caller, "delivery/")[1]
	}
	if strings.Contains(caller, "service/") {
		caller = strings.Split(caller, "service/")[1]
	}
	// if strings.Contains(caller, "repositories/") {
	// 	caller = strings.Split(caller, "repositories/")[1]
	// }
	return caller
}

func shortFuncNameCaller(funcName string) string {
	if strings.Contains(funcName, ServiceName+"/") {
		funcName = strings.Split(funcName, ServiceName+"/")[1]
	}
	if strings.Contains(funcName, "(") {
		funcName = "(" + strings.Split(funcName, "(")[1]
	}
	if strings.Contains(funcName, ".func1") {
		funcName = strings.Split(funcName, ".func1")[0]
	}
	return funcName
}

type CustomLogger struct {
	logger.Config
}

func (c *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *c
	newLogger.LogLevel = level
	return &newLogger
}

func (c *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Info {
		log.Printf("[INFO] "+msg, data...)
	}
}

func (c *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Warn {
		log.Printf("[WARN] "+msg, data...)
	}
}

func (c *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if c.LogLevel >= logger.Error {
		log.Printf("[ERROR] "+msg, data...)
	}
}

func (c *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if c.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	temp := fmt.Sprintf("[%.3fms] [rows:%d] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	switch {
	case err != nil && c.LogLevel >= logger.Error:
		Log.InfoQuery("DEBUG_QUERY_ERROR", zap.Any("query", temp), zap.Error(err))
	case elapsed > c.SlowThreshold && c.SlowThreshold != 0 && c.LogLevel >= logger.Warn:
		Log.InfoQuery("DEBUG_QUERY_SLOW", zap.Any("query", temp))
	case c.LogLevel >= logger.Info:
		Log.InfoQuery("DEBUG_QUERY", zap.Any("query", temp))
	}
}

func ToByte(a interface{}) []byte {
	buffers := new(bytes.Buffer)
	json.NewEncoder(buffers).Encode(a)
	return buffers.Bytes()
}

func GetTimeUTC7() time.Time {
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return now.In(loc)
}

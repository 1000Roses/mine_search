package settings

import (
	"database/sql"
	"errors"
	"fmt"
	"mine/internal"
	"mine/internal/models"
	"mine/internal/utils"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis"
	"github.com/shettyh/threadpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Configs struct {
	Host          string `mapstructure:"DB_HOST"`
	User          string `mapstructure:"DB_USERNAME"`
	Password      string `mapstructure:"DB_PASSWORD"`
	DBName        string `mapstructure:"DB_DATABASE"`
	Port          int    `mapstructure:"DB_PORT"`
	UseProduction bool   `mapstructure:"USE_PRODUCTION"`
	IsDev         bool   `mapstructure:"IS_DEV"`
	TypeSenseKey  string `mapstructure:"TYPESENSE_KEY"`
}
type DateTimeLayout struct {
	YMD     string
	YMD_HMS string
	YDM     string
}
type AppKeys struct {
	AuthClientKey string
}
type AppSettings struct {
	Cfgs              *Configs
	DateFm            *DateTimeLayout
	Keys              *AppKeys
	Accounts          *Accounts
	Regexs            *Regexs
	Log               internal.TempLog
	Logger            *zap.Logger
	ErrMsgs           *ErrMsgs
	Brokers           []string
	KafkaTopicName    string
	KafkaTopicNameAll string
}

func NewSQLDB(config *Configs) *gorm.DB {
	fmt.Println("NewSQLDB start...")
	DBDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
		config.User, config.Password, config.Host, config.Port, config.DBName)
	DBDSN = DBDSN + "&loc=Asia%2FHo_Chi_Minh"
	DB, err := sql.Open("mysql", DBDSN)
	if err != nil {
		// cm.Log.Debug("Open mysql connection failed", zap.Error(err))
		fmt.Println("NewSQLDB err", err)
		return nil
	}
	DB.SetConnMaxLifetime(0)
	DB.SetMaxOpenConns(0)
	DB.SetMaxIdleConns(0)
	errPing := DB.Ping()
	if errPing != nil {
		// cm.Log.Debug("go-mysql-driver ping failed", zap.Error(errPing))
		fmt.Println("DB PING ERR: ", errPing)
		return nil
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: DB,
	}), &gorm.Config{
		Logger: &internal.CustomLogger{
			Config: logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Silent,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		},
	})
	// gormDB, err := gorm.Open(mysql.Dialector{
	// 	Config: &mysql.Config{
	// 		Conn: DB,
	// 	},
	// }, &gorm.Config{
	// 	PrepareStmt: false,
	// })
	if err != nil {
		fmt.Println("gorm DB connect failed", err)
		return nil
	}
	fmt.Println("NewSQLDB success")
	return gormDB
}

func NewAppConfigs() (configs *Configs, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		configs := &Configs{}
		fmt.Println("Load config file error", err)
		db_host, check := utils.GetDefaultEnv("DB_HOST", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		db_port, check := utils.GetDefaultEnv("DB_PORT", "3306")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		db_user, check := utils.GetDefaultEnv("DB_USERNAME", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		db_password, check := utils.GetDefaultEnv("DB_PASSWORD", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		db_name, check := utils.GetDefaultEnv("DB_DATABASE", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		UseProduction, check := utils.GetDefaultEnv("USE_PRODUCTION", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		typesenseKey, check := utils.GetDefaultEnv("TYPESENSE_KEY", "")
		if !check {
			return nil, errors.New("Config missing variable")
		}
		IsDev, check := utils.GetDefaultEnv("IS_DEV", "")
		if !check {
			IsDev = "0"
		}
		port, err_ := strconv.Atoi(db_port)
		if err_ != nil {
			fmt.Println("db_port strconv.Atoi:", db_port, err_.Error())
			port = 3306
		}
		use_product, err_ := strconv.Atoi(UseProduction)
		if err_ != nil {
			fmt.Println("UseProduction strconv.Atoi:", UseProduction, err_.Error())
			use_product = 0
		}
		is_dev, err_ := strconv.Atoi(IsDev)
		if err_ != nil {
			fmt.Println("IsDev strconv.Atoi:", IsDev, err_.Error())
			is_dev = 0
		}
		configs.Host = db_host
		configs.Port = port
		configs.DBName = db_name
		configs.User = db_user
		configs.Password = db_password
		configs.TypeSenseKey = typesenseKey
		if use_product == 0 {
			configs.UseProduction = false
		} else {
			configs.UseProduction = true
		}
		if is_dev == 0 {
			configs.IsDev = false
		} else {
			configs.IsDev = true
		}
		return configs, nil
	} else {
		fmt.Println("Load config file successfully")
		err = viper.Unmarshal(&configs)
		if err != nil {
			fmt.Println("viper.Unmarshal ", err)
			return
		}
		return
	}
}

func NewDateTimeLayout() *DateTimeLayout {
	return &DateTimeLayout{
		YMD:     "2006-01-02",
		YMD_HMS: "2006-01-02 15:04:05",
		YDM:     "2006-02-01",
	}
}

func NewAppKeys(useProduction bool) *AppKeys {
	result := &AppKeys{
		AuthClientKey: "Give me the keys",
	}
	if useProduction {
	}
	return result
}

func NewBrokers(useProduction bool, isDev bool) []string {
	if isDev {
		return []string{"kafka-1:19092", "kafka-2:29092", "kafka-3:39092"}
	}
	if useProduction {
		return []string{"isc-kafka01:9092", "isc-kafka02:9092", "isc-kafka03:9092"}
	}
	return []string{"isc-kafka01:9092", "isc-kafka02:9092", "isc-kafka03:9092"}
}

func NewKafkaTopicNameAll(useProduction bool, isDev bool) string {
	if isDev {
		return "dev-mine"
	}
	if useProduction {
		return "mine"
	}
	return "stag-mine"
}

func NewAppSettings() *AppSettings {
	fmt.Println("NewAppSettings start...")
	appConfigs, err := NewAppConfigs()
	if err != nil {
		fmt.Printf("NewAppSettings err: %v\n", err)
		return nil
	}
	models.Validate = validator.New()
	utils.Pool = threadpool.NewThreadPool(50, 100000)
	fmt.Println("NewAppSettings SUCCESS")
	return &AppSettings{
		Cfgs:   appConfigs,
		DateFm: NewDateTimeLayout(),
		// AppUrls:           NewAppUrls(appConfigs.UseProduction),
		// Endp:              NewAppEndpoints(appConfigs.UseProduction),
		// Ports:             NewServicesPorts(),
		Keys:              NewAppKeys(appConfigs.UseProduction),
		Regexs:            NewRegexs(),
		Log:               internal.TempLog{},
		Logger:            NewLogger(),
		ErrMsgs:           NewErrMsgs(),
		Brokers:           NewBrokers(appConfigs.UseProduction, appConfigs.IsDev),
		KafkaTopicNameAll: NewKafkaTopicNameAll(appConfigs.UseProduction, appConfigs.IsDev),
	}
}

func NewConnectRedis(addr, password string, db int) *redis.Client {
	fmt.Println("NewConnectRedis start...")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	fmt.Println("NewConnectRedis SUCCESS")
	return client
}

type Accounts struct {
	HRLoginUsername string
	HRLoginPassword string
}

type Regexs struct {
	RegexVNPhoneNumber     string
	RegexVNPhoneNumberHide string
	RegexMail              string
	ContractNo             string
}

func NewRegexs() *Regexs {
	result := &Regexs{
		RegexVNPhoneNumber:     `^(84|0)\d{9}$`,
		RegexVNPhoneNumberHide: `^(84|0)\d{6}\*\*\*$`,
		RegexMail:              `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		ContractNo:             `^[a-zA-Z0-9]{1,}$`,
	}
	return result
}

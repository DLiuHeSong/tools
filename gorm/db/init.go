package db

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DefaultConnUrl = "charset=utf8mb4&parseTime=True&loc=Local"

	// defaultConnMaxLifetime, 200 seconds
	defaultConnMaxLifetime = 200

	// defaultMaxIdleConn max idle connects
	defaultMaxIdleConns = 20
)

// DBConfig
type Config struct {
	HostAddress     string // mysql - host:port mysql 地址
	User            string // mysql - user mysql 用户
	Pwd             string // mysql - pwd mysql 密码
	DBName          string // mysql - db mysql 选用 db
	IsLogging       bool   // is log 是否记录日志
	MaxConns        int    // mysql max connections mysql 最大连接数
	MaxIdleConns    int    // mysql max connections 最大空闲连接数
	ConnMaxLifetime int    // ConnMaxLifetime sets the maximum amount of time a connection may be reused
}

type OrmRepository struct {
	C  *Config
	DB *gorm.DB
}

// Init  DB
func (cfg *Config) Init() (*OrmRepository, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.User, cfg.Pwd, cfg.HostAddress, cfg.DBName, DefaultConnUrl)

	loggerLevel := logger.Error
	if cfg.IsLogging {
		loggerLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: loggerLevel, // Log level
				Colorful: true,        // Disable color
			},
		)})
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置最大连接数
	d.SetMaxOpenConns(cfg.MaxConns)

	maxIdleConns := defaultMaxIdleConns
	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}
	// 设置最大空闲连接数
	d.SetMaxIdleConns(maxIdleConns)

	// 设置小于服务器的wait_timeout即可
	connMaxLifetime := defaultConnMaxLifetime
	if cfg.ConnMaxLifetime > 0 {
		connMaxLifetime = cfg.ConnMaxLifetime
	}
	d.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return &OrmRepository{
		C:  cfg,
		DB: db,
	}, nil
}

package init

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
	// DefaultDSNParameters is default DSN paramters combined in string.
	//  - charset: utf8mb4
	//  - parseTime: True
	//  - locale: Local
	DefaultDSNParameters = "charset=utf8mb4&parseTime=True&loc=Local"

	// defaultConnMaxLifetime, 200 seconds
	defaultConnMaxLifetime = 200

	// defaultMaxIdleConn max idle connects
	defaultMaxIdleConns = 20
)

// Config is a config
type Config struct {
	HostAddress     string // mysql - host:port mysql 地址
	Username        string // mysql - user mysql 用户
	Password        string // mysql - pwd mysql 密码
	Database        string // mysql - db mysql 选用 db
	EnableLogging   bool   // is log 是否记录日志
	MaxConnections  int    // mysql max connections mysql 最大连接数
	MaxIdleConns    int    // mysql max connections 最大空闲连接数
	ConnMaxLifetime int    // ConnMaxLifetime sets the maximum amount of time a connection may be reused
}

// CreateRepository creates a Repository
func CreateRepository(cfg *Config) (*gorm.DB, error) {
	// MySQL Connection String
	//   eg: `user:password@tcp(host:3306)/feedcoin?charset=utf8&parseTime=True&loc=utc`
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.Username, cfg.Password, cfg.HostAddress, cfg.Database, DefaultDSNParameters)

	loggerLevel := logger.Error
	if cfg.EnableLogging {
		loggerLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel: loggerLevel, // Log level
				Colorful: true,        // Disable color
			},
		)})
	if err != nil {
		return nil, err
	}

	// Apply gorm settings
	//db.SingularTable(true)

	d, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than MaxIdleConns, then MaxIdleConns will be reduced to match the new MaxOpenConns limit
	// If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
	//
	// https://godoc.org/database/sql#DB.SetMaxOpenConns
	//
	d.SetMaxOpenConns(cfg.MaxConnections)

	// SetMaxIdleConns sets the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
	//
	// If n <= 0, no idle connections are retained.
	maxIdleConns := defaultMaxIdleConns
	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}
	d.SetMaxIdleConns(maxIdleConns)

	// ConnMaxLifetime sets the maximum amount of time a connection may be reused
	// 设置小于服务器的wait_timeout即可
	connMaxLifetime := defaultConnMaxLifetime
	if cfg.ConnMaxLifetime > 0 {
		connMaxLifetime = cfg.ConnMaxLifetime
	}
	d.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	return db, nil
}

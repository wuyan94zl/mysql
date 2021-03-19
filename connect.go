package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB // 定义 mysql 连接实例
var errDb error

type Config struct {
	Host           string
	Port           uint32
	Username       string
	Password       string
	Database       string
	Charset        string
	MaxConnect     int
	MaxIdleConnect int
	MaxLifeSeconds int
}

func ConMysql(config Config) {
	setDefaultVal(&config)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Username, config.Password, config.Host, config.Port, config.Database, config.Charset, true, "Local")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 500 * time.Millisecond, // 慢 SQL 阈值 500ms （如打印所有sql设置1ms）
			LogLevel:      logger.Warn,            // Log level
			Colorful:      false,                  // 禁用彩色打印
		},
	)
	DB, errDb = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if config.MaxConnect > 0 && config.MaxIdleConnect > 0 && config.MaxLifeSeconds > 0{
		sqlDB, _ := DB.DB()
		sqlDB.SetMaxOpenConns(config.MaxConnect)     // 设置最大连接数
		sqlDB.SetMaxIdleConns(config.MaxIdleConnect) //
		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifeSeconds))
	}

	if errDb != nil {
		panic(errDb)
	}
}

func setDefaultVal(config *Config) {
	if config.Host == "" {
		config.Host = "127.0.0.1"
	}
	if config.Port == 0 {
		config.Port = 3306
	}
	if config.Charset == "" {
		config.Charset = "utf8mb4"
	}
}

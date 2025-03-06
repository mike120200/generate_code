package pgconn

import (
	"binbin/common/config"
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// DbInit 初始化数据库
func DbInit() error {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/pg_conn/pg_conn.go DbInit")

	// 从配置中获取数据库参数
	host, err := config.GetConfig("PostgresDB.host")
	if err != nil {
		logger.Error("get postgresDB.host failed:", zap.Error(err))
		return err
	}
	port, err := config.GetConfig("PostgresDB.port")
	if err != nil {
		logger.Error("get postgresDB.port failed:", zap.Error(err))
		return err
	}
	userName, err := config.GetConfig("PostgresDB.user")
	if err != nil {
		logger.Error("get postgresDB.user failed:", zap.Error(err))
		return err
	}
	password, err := config.GetConfig("PostgresDB.password")
	if err != nil {
		logger.Error("get postgresDB.password failed:", zap.Error(err))
		return err
	}
	database, err := config.GetConfig("PostgresDB.database")
	if err != nil {
		logger.Error("get postgresDB.database failed:", zap.Error(err))
		return err
	}

	// 构造 DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, userName, password, database, port)

	// 初始化 GORM
	var errDB error
	db, errDB = gorm.Open(postgres.Open(dsn))
	if errDB != nil {
		logger.Error("failed to connect to database", zap.Error(errDB))
		return errDB
	}

	// 获取底层数据库连接池
	sqlDB, errDB := db.DB()
	if errDB != nil {
		logger.Error("failed to get underlying database connection", zap.Error(errDB))
		return errDB
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(8)                  // 最小连接数
	sqlDB.SetMaxOpenConns(28)                 // 最大连接数
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) // 最大空闲时间
	sqlDB.SetConnMaxLifetime(1 * time.Hour)   // 最大存活时间

	// 测试连接
	errDB = sqlDB.PingContext(context.Background())
	if errDB != nil {
		logger.Error("failed to ping database", zap.Error(errDB))
		return errDB
	}

	logger.Info("database connection initialized successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() (*gorm.DB, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/pg_conn/pg_conn.go GetDB")

	if db == nil {
		return nil, fmt.Errorf("db is nil while it shouldn't")
	}
	return db, nil
}

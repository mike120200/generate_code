package pgconn

import (
	"context"
	"fmt"
	config "test_binbin/common/be_config"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var dbpool *pgxpool.Pool

// DbInit 初始化数据库
func DbInit() error {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/pg_conn/pg_conn.go DbInit")
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
	UserName, err := config.GetConfig("PostgresDB.user")
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
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", UserName, password, host, port, database)
	//创建连接池配置
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Error("parse config failed:", zap.Error(err))
		return err
	}
	//最大连接给到16
	config.MaxConns = 16
	//最小连接数
	config.MinConns = 8
	//最大空闲时间
	config.MaxConnIdleTime = 5 * time.Minute
	//最久存活时间
	config.MaxConnLifetime = 1 * time.Hour

	deadline := time.Now().Add(1 * time.Second)
	cause := fmt.Errorf("pool connection timed out")
	ctx, cancel := context.WithDeadlineCause(context.Background(), deadline, cause)
	defer cancel()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error("new pool failed:", zap.Error(err))
		return err
	}
	dbpool = pool
	// 测试数据库连接池是否有效
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		logger.Error("连接出现问题:", zap.Error(err))
		return err
	}
	//释放连接
	conn.Release()

	return nil
}

// GetPool 获取数据库
func GetPool() (*pgxpool.Pool, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/pg_conn/pg_conn.go GetPool")
	if dbpool == nil {
		return nil, fmt.Errorf("dbpool is nil while it's shouldn't")
	}
	// 测试数据库连接池是否有效
	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		logger.Error("连接出现问题:", zap.Error(err))
		return nil, fmt.Errorf("failed to acquire")
	}
	//释放连接
	conn.Release()
	logger.Info("get pool successfully ")
	return dbpool, nil
}

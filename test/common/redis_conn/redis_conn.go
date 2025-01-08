package redisconn

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	config "test_binbin/common/be_config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var redisDb *redis.Client

func Redis_init() error {
	logger := zap.L()
	if logger == nil {
		fmt.Println(" create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/redis_conn.redis_init")
	host, err := config.GetConfig("Redis.host")
	if err != nil {
		logger.Error("get config Redis.host failed error: " + err.Error())
		return err
	}
	port, err := config.GetConfig("Redis.port")
	if err != nil {
		logger.Error("get config Redis.port failed error: " + err.Error())
		return err
	}
	address := fmt.Sprintf("%s:%s", host, port)
	password, err := config.GetConfig("Redis.password")
	if err != nil {
		logger.Error("get config Redis.password failed error: " + err.Error())
		return err
	}

	DB, err := config.GetConfig("Redis.DB")
	if err != nil {
		logger.Error("get config Redis.DB failed error: " + err.Error())
		return err
	}
	intDB, err := strconv.Atoi(DB)
	if err != nil {
		logger.Error("DB change string to int failed error: " + err.Error())
		return err
	}

	MaxIdleConns, err := config.GetConfig("Redis.MaxIdleConns")
	if err != nil {
		logger.Error("get config Redis.MaxIdleConns failed error: " + err.Error())
		return err
	}
	intMaxIdleConns, err := strconv.Atoi(MaxIdleConns)
	if err != nil {
		logger.Error("MaxIdleConns change string to int failed error: " + err.Error())
		return err
	}

	MaxActiveConns, err := config.GetConfig("Redis.MaxActiveConns")
	if err != nil {
		logger.Error("get config Redis.MaxActiveConns failed error: " + err.Error())
		return err
	}
	IntMaxActiveConns, err := strconv.Atoi(MaxActiveConns)
	if err != nil {
		logger.Error("MaxActiveConns change string to int failed error: " + err.Error())
		return err
	}
	client := redis.NewClient(&redis.Options{
		Addr:           address,
		Password:       password,
		DB:             intDB,
		MaxIdleConns:   intMaxIdleConns,
		MaxActiveConns: IntMaxActiveConns,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//测试
	_, err = client.Ping(ctx).Result()
	if err != nil {
		logger.Sugar().Errorf("redis ping error: %v", err)
		return err
	}
	logger.Info("redis init success")
	redisDb = client
	return nil
}

func GetRedisClient() (*redis.Client, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println("create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/pg_conn/redis_conn.go GetRedisClient")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//测试
	_, err := redisDb.Ping(ctx).Result()
	if err != nil {
		logger.Sugar().Errorf("redis ping error: %v", err)
		return nil, err
	}
	logger.Info("redis init success")
	return redisDb, nil
}

package generated_code

import "strings"

type GeneratedCode string

// projectName是项目的名称，需要替换
func (code GeneratedCode) ReplaceProjectName(prjName string) string {
	return strings.Replace(string(code), "projectName", prjName, -1)
}

// 需要生成的代码中包含项目名称的部分，需要用“projectName”替换
// 例如："projectName/common/config"
var LogCode GeneratedCode = `package zap_log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logDir  = "./logs/"
	logFile = "log.txt"
)

func LoggerInit() error {

	// 检查目录是否存在
	if _, err := os.Stat(logDir); !os.IsNotExist(err) {
		// 目录已存在
		fmt.Println("目录已存在：")
	} else {
		// 目录不存在，创建目录
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Printf("无法创建目录：%v\n", err)
			return nil
		}
		fmt.Println("目录已创建：")
	}
	fileConfig := &lumberjack.Logger{
		Filename:   logDir + logFile, //日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    5,                //文件大小限制,单位MB
		MaxBackups: 5,                //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,            //是否压缩处理
	}
	//zapcore.AddSync 函数通常用于将一个 io.Writer 转换为 zapcore.WriteSyncer 接口的实现。
	//zapcore.WriteSyncer 是 zap 包中的一个接口，它扩展了 io.Writer 接口，增加了一个 Sync 方法，
	//该方法用于确保所有已写入的数据都被正确地刷新到它们的最终目的地
	FileWriteSyncer := zapcore.AddSync(fileConfig)
	stdioWriteSyncer := zapcore.AddSync(os.Stdout)

	//设置日志编码器
	EncoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(EncoderConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, FileWriteSyncer, zap.InfoLevel),
		zapcore.NewCore(encoder, stdioWriteSyncer, zap.DebugLevel),
	)
	//初始化实例
	logger := zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(logger)
	logger.Info("log init success")
	return nil
}
`
var LogCodeTest GeneratedCode = `package zap_log_test

import (
	"fmt"
	zaplog "projectName/common/log"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {

	logger := zap.L()
	if logger == nil {
		fmt.Println(" create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("success")
}

func init() {
	if err := zaplog.LoggerInit(); err != nil {
		panic(err)
	}
}
`
var PgconnCode GeneratedCode = `package pgconn

import (
	"context"
	"fmt"
	config "projectName/common/config"
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
`
var PgconnCodeTest GeneratedCode = `package pgconn_test

import (
	"context"
	config "projectName/common/config"
	Log "projectName/common/log"
	pgconn "projectName/common/pg_conn"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPool(t *testing.T) {
	should := assert.New(t)
	pool, err := pgconn.GetPool()
	if should.NoError(err) {
		conn, err := pool.Acquire(context.Background())
		if should.NoError(err) {
			err := conn.Ping(context.Background())
			if should.NoError(err) {
				defer conn.Release()
				t.Log("ping success")
			}

		}
	}
}

func init() {
	if err := Log.LoggerInit(); err != nil {
		panic(err)
	}
	if err := config.ViperInit(1); err != nil {
		panic(err)
	}
	if err := pgconn.DbInit(); err != nil {
		panic(err)
	}
}
`
var PgconnCode_gorm GeneratedCode = `package pgconn

import (
	"projectName/common/config"
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
`
var PgconnCodeTest_gorm GeneratedCode = `package pgconn_test

import (
	config "projectName/common/config"
	Log "projectName/common/log"
	pgconn "projectName/common/pg_conn"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPool(t *testing.T) {
	should := assert.New(t)
	db, err := pgconn.GetDB()
	if should.NoError(err) {
		// 获取底层数据库连接池
		sqlDB, errDB := db.DB()
		if errDB != nil {
			t.Error("failed to get underlying database connection" + errDB.Error())
			return
		}
		//测试连接池是否正常
		err = sqlDB.PingContext(context.Background())
		if !should.NoError(err) {
			t.Error("failed to ping database" + err.Error())
			return
		}
	}
}

func init() {
	if err := Log.LoggerInit(); err != nil {
		panic(err)
	}
	if err := config.ViperInit(1); err != nil {
		panic(err)
	}
	if err := pgconn.DbInit(); err != nil {
		panic(err)
	}
}
`

var RedisCode GeneratedCode = `package redisconn

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
	config "projectName/common/config"

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
`
var RedisCodeTest GeneratedCode = `package redisconn_test

import (
	"context"
	config "projectName/common/config"
	Log "projectName/common/log"
	redis_conn "projectName/common/redis_conn"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRedisClient(t *testing.T) {
	should := assert.New(t)
	pool, err := redis_conn.GetRedisClient()
	if should.NoError(err) {
		conn := pool.Ping(context.Background())
		t.Log(conn.Result())
	}
}

func init() {
	if err := Log.LoggerInit(); err != nil {
		panic(err)
	}
	if err := config.ViperInit(1); err != nil {
		panic(err)
	}
	if err := redis_conn.Redis_init(); err != nil {
		panic(err)
	}
}
`
var ConfigCode GeneratedCode = `package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ViperInit 初始化viper
func ViperInit(mode int) error {
	if mode == 1 {
		viper.SetConfigName(".conf_linux_env")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath("../..")
		viper.AddConfigPath("../../..")
		viper.AddConfigPath("../../../..")
		viper.AddConfigPath("../../../../..")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return nil
	} else {
		viper.SetConfigName(".conf_linux")
		viper.SetConfigType("json")
		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath("../..")
		viper.AddConfigPath("../../..")
		viper.AddConfigPath("../../../..")
		viper.AddConfigPath("../../../../..")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		return nil
	}
}

// GetConfig 获取配置
func GetConfig(key string) (string, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println(" create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/config/config.go GetConfig")
	if key == "" {
		logger.Error("the key is \"\" ")
		return "", fmt.Errorf("the key is \"\" ")
	}
	if !viper.IsSet(key) {
		logger.Error("the key " + key + " does not exist")
		return "", fmt.Errorf("the key " + key + " does not exist")
	}
	config := viper.GetString(key)
	if config == "" {
		logger.Error("the value of key " + key + " is empty")
		return "", fmt.Errorf("the value of key " + key + " is empty")
	}
	logger.Info("get config success, key: " + key + ", value: " + config)
	return config, nil
}


`
var ConfigCodeTest GeneratedCode = `package config_test

import (
	config "projectName/common/config"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViper(t *testing.T) {
	should := assert.New(t)
	result, err := config.GetConfig("Redis.DB")
	if should.NoError(err) {
		t.Log(result)
		intResult, err := strconv.Atoi(result)
		if should.NoError(err) {
			t.Log(intResult)
		}
	}

}

func init() {
	if err := config.ViperInit(1); err != nil {
		panic(err)
	}

}
`


var ResultGeneratedCode GeneratedCode= `package result

import (
	"encoding/json"
	"errors"
)

var (
	ErrMsgEmpty = errors.New("response msg is empty")
)

// 定义响应码
const (
	SuccessCode = 200
)

// Response 统一响应结构体
type Response struct {
	Code int         ` + "`json:\"code\"`" + `
	Msg  string      ` + "`json:\"msg\"`" + `
	Data interface{} ` + "`json:\"data\"`" + `
}

// NewResponse 创建响应对象
// code 响应码
// msg 响应信息
// data 响应数据，如果发生错误的话，这里就为空
func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// ToJson 将响应对象转换为json格式
func (response *Response) ToJson() ([]byte, error) {
	if response.Msg == "" {
		return nil, ErrMsgEmpty
	}
	// 将结构体转换为json
	jsonData, err := json.Marshal(response)
	if err != nil {
		return nil, errors.New("json.Marshal failed: " + err.Error())
	}
	return jsonData, nil
}
`

var ResultCodeTest GeneratedCode= `package result_test

import (
	"projectName/common/result"
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	var cases = []struct {
		test_name string
		msg       string
		code      int
		data      interface{}
	}{
		{"common_test", "successfully", 200, "hello"},
		{"empty_msg_test", "", 699, nil},
	}
	var expected = []struct {
		Err  error
		Msg  string
		Code int
		Data interface{}
	}{
		{nil, "successfully", 200, "hello"},
		{result.ErrMsgEmpty, "", 200, nil},
	}
	for _, c := range cases {
		t.Run(c.test_name, func(t *testing.T) {
			response := result.NewResponse(c.code, c.msg, c.data)
			if reflect.DeepEqual(response, expected) {
				t.Errorf("expected: %v, got: %v", expected, response)
			}
		})
	}
}
func TestToJson(t *testing.T) {
	var cases = []struct {
		test_name string
		msg       string
		code      int
		data      interface{}
	}{
		{"common_test", "successfully", 200, "hello"},
	}
	for _, c := range cases {
		t.Run(c.test_name, func(t *testing.T) {
			response := result.NewResponse(c.code, c.msg, c.data)
			jsonData, err := response.ToJson()
			if err != nil {
				t.Errorf("json marshal failed: " + err.Error())
				return
			}
			t.Logf("%v", string(jsonData))
		})
	}
}
`
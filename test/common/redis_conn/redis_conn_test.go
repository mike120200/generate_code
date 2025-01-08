package redisconn_test

import (
	"context"
	config "test_binbin/common/be_config"
	Log "test_binbin/common/log"
	redis_conn "test_binbin/common/redis_conn"
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
	if err := config.ViperInit(); err != nil {
		panic(err)
	}
	if err := redis_conn.Redis_init(); err != nil {
		panic(err)
	}
}

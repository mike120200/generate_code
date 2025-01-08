package pgconn_test

import (
	"context"
	config "test_binbin/common/be_config"
	Log "test_binbin/common/log"
	pgconn "test_binbin/common/pg_conn"
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
	if err := config.ViperInit(); err != nil {
		panic(err)
	}
	if err := pgconn.DbInit(); err != nil {
		panic(err)
	}
}

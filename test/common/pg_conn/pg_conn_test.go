package pgconn_test

import (
	config "binbin/common/config"
	Log "binbin/common/log"
	pgconn "binbin/common/pg_conn"
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

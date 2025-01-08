package config_test

import (
	be_config "test_binbin/common/be_config"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViper(t *testing.T) {
	should := assert.New(t)
	result, err := be_config.GetConfig("Redis.DB")
	if should.NoError(err) {
		t.Log(result)
		intResult, err := strconv.Atoi(result)
		if should.NoError(err) {
			t.Log(intResult)
		}
	}

}

func init() {
	if err := be_config.ViperInit(); err != nil {
		panic(err)
	}

}

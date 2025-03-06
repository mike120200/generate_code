package config_test

import (
	config "binbin/common/config"
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

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ViperInit 初始化viper
func ViperInit() error {
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

// GetConfig 获取配置
func GetConfig(key string) (string, error) {
	logger := zap.L()
	if logger == nil {
		fmt.Println(" create logger failed, please check zap logger")
		os.Exit(-1)
	}
	logger.Info("-->go/common/be_config/be_config.go GetConfig")
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


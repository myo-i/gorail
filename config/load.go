package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbPath string
}

// 環境変数を取得
func Load() (Config, error) {
	viper.SetConfigName("db")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to ")
	}

	return config, err
}

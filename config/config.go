package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type LogConfig struct {
	Use    bool
	Logger string
}

type ORMConfig struct {
	Use   bool
	Frame string
}

type WebConfig struct {
	Use   bool
	Frame string
}

type RedisConfig struct {
	Use bool
}

var Info struct {
	Project struct {
		Name   string
		Module string
	}
	Log   LogConfig
	ORM   ORMConfig
	Web   WebConfig
	Redis RedisConfig
}

func Setup(filePath string) {
	viper.SetConfigType("toml")

	viper.SetConfigFile(filePath)
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	if err := viper.ReadInConfig(); err != nil {
		panic("Load config failed!")
	}

	viper.SetDefault("log", LogConfig{Use: true, Logger: "go.uber.org/zap"})
	viper.SetDefault("orm", ORMConfig{Use: true, Frame: "gorm.io/gorm"})
	viper.SetDefault("web", WebConfig{Use: true, Frame: "github.com/gin-gonic/gin"})
	viper.SetDefault("redis", RedisConfig{Use: true})

	err := viper.Unmarshal(&Info)
	if err != nil {
		panic("Marshal config failed!")
	}
}

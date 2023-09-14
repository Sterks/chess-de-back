package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Http     HttpConfig `mapstructure:"http"`
	Database Database   `mapstructure:"database"`
	MongoDB  MongoDB    `mapstructure:"mongo_db"`
}

type HttpConfig struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	ReadTimeout    time.Duration `mapstructure:"readTimeout"`
	WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Port     string `mapstructure:"port"`
	SslMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`
}

type MongoDB struct {
	URL string `mapstructure:"url"`
}

func InitConfig(configPath string) (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("main")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg *Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

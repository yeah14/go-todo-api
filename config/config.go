package config

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Server   SeverConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Log      LogConfig
}

type SeverConfig struct {
	Port         string
	Mode         string
	ReadTimeout  int
	WriteTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

type JWTConfig struct {
	Secret        string
	AccessExpire  int
	RefreshExpire int
	issuer        string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type LogConfig struct {
	Level      string
	File       string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

var GlobalConfig config

func InitConfig(configFile string) {
	if configFile == "" {
		configFile = "config/dev.yaml"
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取配置文件错误：%v", err)
	}
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatal("解析配置文件错误:%v", err)
	}
	//fmt.Println(GlobalConfig.Database.Host, GlobalConfig.Database.Port, GlobalConfig.Database.Username, GlobalConfig.Database.Password)
}

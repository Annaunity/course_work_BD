package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"coursework/pkg/log"
	"sync"
)

type Config struct {
	LoggerLevel string `yaml:"logger_level"`
	Listen      struct {
		BindIp string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"6000"`
	} `yaml:"listen"`
	JWT struct {
		Secret string `yaml:"secret" env-required:"true"`
	} `yaml:"jwt" env-required:"true"`
	Postgres struct {
		Host     string `yaml:"host" env-required:"true"`
		Port     string `yaml:"port" env-required:"true"`
		Username string `yaml:"username" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		Database string `yaml:"database" env-required:"true"`
	} `yaml:"postgres" env-default:"true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := log.GetLogger()
		logger.Println("try to open and read config file")

		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Println(help)
			logger.Fatalln(err)
		}
	})

	return instance
}
